package vps

import (
	"context"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"

	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("vps")

// Manager manager order
type Manager struct {
	*db.SQLDB
	cfg       config.MallCfg
	vpsClient map[string]*ecs20140526.Client

	getInstanceInfoRunning bool
}

// NewManager returns a new manager instance
func NewManager(sdb *db.SQLDB, getCfg dtypes.GetMallConfigFunc) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		SQLDB:     sdb,
		cfg:       cfg,
		vpsClient: make(map[string]*ecs20140526.Client),
	}
	go m.cronGetInstanceDefaultInfo()

	return m, nil
}

// CreateAliYunInstance creates an Alibaba Cloud instance.
func (m *Manager) CreateAliYunInstance(vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret

	priceUnit := vpsInfo.PeriodUnit
	period := vpsInfo.Period
	regionID := vpsInfo.RegionID
	if priceUnit == "Year" {
		priceUnit = "Month"
		period = period * 12
	}

	var securityGroupID string

	securityGroups, err := aliyun.DescribeSecurityGroups(regionID, accessKeyID, accessKeySecret)
	if err == nil && len(securityGroups) > 0 {
		securityGroupID = securityGroups[0]
	} else {
		securityGroupID, err = aliyun.CreateSecurityGroup(regionID, accessKeyID, accessKeySecret)
		if err != nil {
			log.Errorf("CreateSecurityGroup err: %s", err.Error())
			return nil, xerrors.New(err.Error())
		}
	}

	log.Debugln("securityGroupID:", securityGroupID, " , DryRun:", vpsInfo.DryRun)

	result, err := aliyun.CreateInstance(accessKeyID, accessKeySecret, vpsInfo, vpsInfo.DryRun)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(err.Error())
	}

	address, err := aliyun.AllocatePublicIPAddress(regionID, accessKeyID, accessKeySecret, result.InstanceID)
	if err != nil {
		log.Errorf("AllocatePublicIpAddress err: %s", err.Error())
	} else {
		result.PublicIPAddress = address
	}

	// 设置安全端口 (使用账密的时候必须用)
	// err = aliyun.AuthorizeSecurityGroup(regionID, k, s, securityGroupID)
	// if err != nil {
	// 	log.Errorf("AuthorizeSecurityGroup err: %s", err.Error())
	// }

	go func() {
		time.Sleep(1 * time.Minute)

		err = aliyun.StartInstance(regionID, accessKeyID, accessKeySecret, result.InstanceID)
		if err != nil {
			log.Infoln("StartInstance err:", err)
		}
	}()

	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	instanceInfo, err := aliyun.DescribeInstances(regionID, accessKeyID, accessKeySecret, instanceIds)
	if err != nil {
		log.Errorf("DescribeInstances err: %s", err.Error())
	} else {
		if len(instanceInfo.Body.Instances.Instance) > 0 {
			ip := instanceInfo.Body.Instances.Instance[0].PublicIpAddress.IpAddress[0]
			securityGroupId := ""
			if len(instanceInfo.Body.Instances.Instance) > 0 {
				securityGroupId = *instanceInfo.Body.Instances.Instance[0].SecurityGroupIds.SecurityGroupId[0]
			}
			OSType := instanceInfo.Body.Instances.Instance[0].OSType
			InstanceName := instanceInfo.Body.Instances.Instance[0].InstanceName
			BandwidthOut := instanceInfo.Body.Instances.Instance[0].InternetMaxBandwidthOut
			Cores := instanceInfo.Body.Instances.Instance[0].Cpu
			Memory := instanceInfo.Body.Instances.Instance[0].Memory
			ExpiredTime := instanceInfo.Body.Instances.Instance[0].ExpiredTime
			instanceDetailsInfo := &types.InstanceDetails{
				IPAddress:       *ip,
				InstanceID:      result.InstanceID,
				SecurityGroupID: securityGroupId,
				OSType:          *OSType,
				Cores:           *Cores,
				Memory:          float32(*Memory),
				InstanceName:    *InstanceName,
				ExpiredTime:     *ExpiredTime,
				BandwidthOut:    *BandwidthOut,
				AccessKey:       result.AccessKey,
			}
			errU := m.UpdateInstanceInfoOfUser(instanceDetailsInfo)
			if errU != nil {
				log.Errorf("UpdateVpsInstance:%v", errU)
			}
		}
	}

	return result, nil
}

// RenewInstance renews an instance.
func (m *Manager) RenewInstance(renewInstanceRequest *types.RenewInstanceRequest) error {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret

	err := aliyun.RenewInstance(accessKeyID, accessKeySecret, renewInstanceRequest)
	if err != nil {
		log.Errorf("RenewInstance err: %s", err.Error())
		return xerrors.New(err.Error())
	}
	return nil
}

// cronFetchInstanceDefaultInfo fetches default instance information periodically.
func (m *Manager) cronGetInstanceDefaultInfo() {
	now := time.Now()

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
	if now.After(nextTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	duration := nextTime.Sub(now)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		<-timer.C

		log.Debugln("start instance check ")

		timer.Reset(24 * time.Hour)

		m.UpdateInstanceDefaultInfo("")
	}
}

// UpdateInstanceDefaultInfo updates default instance information.
func (m *Manager) UpdateInstanceDefaultInfo(regionID string) {
	if m.getInstanceInfoRunning {
		log.Debugln("task is running")
		return
	}

	m.getInstanceInfoRunning = true
	defer func() {
		m.getInstanceInfoRunning = false
	}()

	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret
	var ctx context.Context

	var regionIDs []string
	if regionID != "" {
		regionIDs = append(regionIDs, regionID)
	} else {
		regions, err := aliyun.DescribeRegions(accessKeyID, accessKeySecret)
		if err != nil {
			log.Errorf("DescribePrice err:%v", err.Error())
			return
		}
		for _, region := range regions.Body.Regions.Region {
			regionIDs = append(regionIDs, *region.RegionId)
		}
	}

	for _, regionID := range regionIDs {
		instanceType := &types.DescribeInstanceTypeReq{
			RegionID:     regionID,
			CPUCoreCount: 0,
			MemorySize:   0,
		}

		instances, err := m.DescribeInstanceType(ctx, instanceType)
		if err != nil {
			log.Errorf("DescribeInstanceType err:%v", err.Error())
			continue
		}
		for _, instance := range instances.InstanceTypes {
			ok, err := m.InstancesDefaultExists(instance.InstanceTypeID, regionID)
			if err != nil {
				log.Errorf("InstancesDefaultExists err:%v", err.Error())
				continue
			}
			if ok {
				continue
			}
			images, err := m.DescribeImages(ctx, regionID, instance.InstanceTypeID)
			if err != nil {
				log.Errorf("DescribePrice err:%v", err.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeID, regionID)
				continue
			}
			disk := &types.AvailableResourceReq{
				InstanceType:        instance.InstanceTypeID,
				RegionID:            regionID,
				DestinationResource: "SystemDisk",
			}

			disks, err := m.DescribeAvailableResourceForDesk(ctx, disk)
			if err != nil {
				log.Errorf("DescribePrice err:%v", err.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeID, regionID)
				continue
			}
			if len(disks) > 0 {
				priceReq := &types.DescribePriceReq{
					RegionID:                regionID,
					InstanceType:            instance.InstanceTypeID,
					PriceUnit:               "Month",
					ImageID:                 images[0].ImageID,
					InternetChargeType:      "PayByTraffic",
					SystemDiskCategory:      disks[0].Value,
					SystemDiskSize:          40,
					Period:                  1,
					Amount:                  1,
					InternetMaxBandwidthOut: 10,
				}
				price, err := aliyun.DescribePrice(accessKeyID, accessKeySecret, priceReq)
				if err != nil {
					log.Errorf("DescribePrice err:%v", err.Error())
					_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeID, regionID)
					continue
				}
				info := &types.DescribeInstanceTypeFromBase{
					RegionID:               regionID,
					InstanceTypeID:         instance.InstanceTypeID,
					MemorySize:             instance.MemorySize,
					CPUArchitecture:        instance.CPUArchitecture,
					InstanceCategory:       instance.InstanceCategory,
					CPUCoreCount:           instance.CPUCoreCount,
					AvailableZone:          instance.AvailableZone,
					InstanceTypeFamily:     instance.InstanceTypeFamily,
					PhysicalProcessorModel: instance.PhysicalProcessorModel,
					OriginalPrice:          price.OriginalPrice,
					Price:                  price.USDPrice,
					Status:                 instance.Status,
				}
				saveErr := m.SaveInstancesInfo(info)
				if err != nil {
					log.Errorf("SaveMyInstancesInfo:%v", saveErr)
				}

			}

		}

	}
	return
}

// DescribeInstanceType fetches instance type information.
func (m *Manager) DescribeInstanceType(ctx context.Context, instanceType *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error) {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret
	rsp, err := aliyun.DescribeInstanceTypes(accessKeyID, accessKeySecret, instanceType)
	if err != nil {
		log.Errorf("DescribeInstanceTypes err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	availableResource, err := aliyun.DescribeAvailableResource(accessKeyID, accessKeySecret, instanceType)
	if err != nil {
		log.Errorf("DescribeAvailableResource err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	if availableResource.Body.AvailableZones == nil {
		return nil, xerrors.New("parameter error")
	}

	rspDataList := &types.DescribeInstanceTypeResponse{
		NextToken: *rsp.Body.NextToken,
	}

	availableZone := len(availableResource.Body.AvailableZones.AvailableZone)
	if availableZone < 0 {
		return rspDataList, nil
	}

	instanceTypes := make(map[string]string)
	for _, data := range availableResource.Body.AvailableZones.AvailableZone {
		availableTypes := data.AvailableResources.AvailableResource
		if len(availableTypes) > 0 {
			for _, instanceTypeResource := range availableTypes {
				Resources := instanceTypeResource.SupportedResources.SupportedResource
				if len(Resources) > 0 {
					for _, Resource := range Resources {
						instanceTypes[*Resource.Value] = *Resource.Status
					}
				}
			}
		}
	}

	for _, data := range rsp.Body.InstanceTypes.InstanceType {
		if v, ok := instanceTypes[*data.InstanceTypeId]; ok {
			rspData := &types.DescribeInstanceType{
				InstanceCategory:       *data.InstanceCategory,
				InstanceTypeID:         *data.InstanceTypeId,
				MemorySize:             *data.MemorySize,
				CPUArchitecture:        *data.CpuArchitecture,
				InstanceTypeFamily:     *data.InstanceTypeFamily,
				CPUCoreCount:           *data.CpuCoreCount,
				AvailableZone:          availableZone,
				PhysicalProcessorModel: *data.PhysicalProcessorModel,
				Status:                 v,
			}
			rspDataList.InstanceTypes = append(rspDataList.InstanceTypes, rspData)
		}
	}

	return rspDataList, nil
}

// DescribeImages fetches image information.
func (m *Manager) DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error) {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret

	rsp, err := aliyun.DescribeImages(regionID, accessKeyID, accessKeySecret, instanceType)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	var rspDataList []*types.DescribeImageResponse
	for _, data := range rsp.Body.Images.Image {
		rspData := &types.DescribeImageResponse{
			ImageID:      *data.ImageId,
			ImageFamily:  *data.ImageFamily,
			ImageName:    *data.ImageName,
			Architecture: *data.Architecture,
			OSName:       *data.OSName,
			OSType:       *data.OSType,
			Platform:     *data.Platform,
		}
		rspDataList = append(rspDataList, rspData)
	}

	return rspDataList, nil
}

// DescribeAvailableResourceForDesk fetches available resources for the system disk.
func (m *Manager) DescribeAvailableResourceForDesk(ctx context.Context, desk *types.AvailableResourceReq) ([]*types.AvailableResourceResponse, error) {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret
	rsp, err := aliyun.DescribeAvailableResourceForDesk(accessKeyID, accessKeySecret, desk)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	category := map[string]int{
		"cloud":            1,
		"cloud_essd":       1,
		"cloud_ssd":        1,
		"cloud_efficiency": 1,
		"ephemeral_ssd":    1,
	}

	var rspDataList []*types.AvailableResourceResponse
	if rsp.Body.AvailableZones == nil {
		return rspDataList, nil
	}

	if len(rsp.Body.AvailableZones.AvailableZone) > 0 {
		availableResources := rsp.Body.AvailableZones.AvailableZone[0].AvailableResources.AvailableResource
		if len(availableResources) > 0 {
			systemDesk := availableResources[0].SupportedResources.SupportedResource
			if len(systemDesk) > 0 {
				for _, data := range systemDesk {
					if *data.Status == "Available" {
						if _, ok := category[*data.Value]; ok {
							desk := &types.AvailableResourceResponse{
								Min:   *data.Min,
								Max:   *data.Max,
								Value: *data.Value,
								Unit:  *data.Unit,
							}
							rspDataList = append(rspDataList, desk)
						}
					}
				}
			}
		}
	}

	return rspDataList, nil
}

// ModifyInstanceRenew modifies instance renewal settings.
func (m *Manager) ModifyInstanceRenew(renewReq *types.SetRenewOrderReq) error {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret
	err := m.UpdateRenewInstanceStatus(renewReq)
	if err != nil {
		log.Errorf("UpdateRenewInstanceStatus:%v", err)
		return err
	}

	errSDK := aliyun.ModifyInstanceAutoRenewAttribute(accessKeyID, accessKeySecret, renewReq)
	if err != nil {
		log.Errorf("ModifyInstanceAutoRenewAttribute err: %s", *errSDK.Message)
		return &api.ErrWeb{Code: terrors.ThisInstanceNotSupportOperation.Int(), Message: *errSDK.Message}
	}

	return nil
}
