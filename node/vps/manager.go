package vps

import (
	"context"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"

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
	cfg config.MallCfg

	getInstanceInfoRunning bool
}

// NewManager returns a new manager instance
func NewManager(sdb *db.SQLDB, getCfg dtypes.GetMallConfigFunc) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		SQLDB: sdb,
		cfg:   cfg,
	}
	go m.cronGetInstanceDefaultInfo()

	return m, nil
}

// CreateAliYunInstance creates an Alibaba Cloud instance.
func (m *Manager) CreateAliYunInstance(orderID string, vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret

	priceUnit := vpsInfo.PeriodUnit
	period := vpsInfo.Period
	regionID := vpsInfo.RegionId
	if priceUnit == "Year" {
		priceUnit = "Month"
		period = period * 12
	}

	var securityGroupID string

	securityGroups, sErr := aliyun.DescribeSecurityGroups(regionID, accessKeyID, accessKeySecret)
	if sErr == nil && len(securityGroups) > 0 {
		securityGroupID = securityGroups[0]
	} else {
		securityGroupID, sErr = aliyun.CreateSecurityGroup(regionID, accessKeyID, accessKeySecret)
		if sErr != nil {
			log.Errorf("CreateSecurityGroup err: %v", sErr)
			return nil, xerrors.New(*sErr.Message)
		}
	}

	log.Debugln("securityGroupID:", securityGroupID, " , DryRun:", vpsInfo.DryRun)

	result, sErr := aliyun.CreateInstance(accessKeyID, accessKeySecret, vpsInfo, vpsInfo.DryRun)
	if sErr != nil {
		log.Errorf("CreateInstance err: %v", sErr)
		return nil, xerrors.New(*sErr.Message)
	}

	address, sErr := aliyun.AllocatePublicIPAddress(regionID, accessKeyID, accessKeySecret, result.InstanceID)
	if sErr != nil {
		log.Errorf("AllocatePublicIpAddress err: %v", sErr)
	} else {
		result.PublicIpAddress = address
	}

	// 设置安全端口 (使用账密的时候必须用)
	// err = aliyun.AuthorizeSecurityGroup(regionID, accessKeyID, accessKeySecret, securityGroupID)
	// if err != nil {
	// 	log.Errorf("AuthorizeSecurityGroup err: %s", err.Error())
	// }

	go func() {
		time.Sleep(1 * time.Minute)

		sErr = aliyun.StartInstance(regionID, accessKeyID, accessKeySecret, result.InstanceID)
		if sErr != nil {
			log.Infoln("StartInstance err:", sErr)
		}
	}()

	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	instanceInfo, sErr := aliyun.DescribeInstances(regionID, accessKeyID, accessKeySecret, instanceIds)
	if sErr != nil {
		log.Errorf("DescribeInstances err: %v", sErr)
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
				IpAddress:       *ip,
				InstanceId:      result.InstanceID,
				SecurityGroupId: securityGroupId,
				OSType:          *OSType,
				Cores:           *Cores,
				Memory:          float32(*Memory),
				InstanceName:    *InstanceName,
				ExpiredTime:     *ExpiredTime,
				BandwidthOut:    *BandwidthOut,
				AccessKey:       result.AccessKey,
				OrderID:         orderID,
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

	sErr := aliyun.RenewInstance(accessKeyID, accessKeySecret, renewInstanceRequest)
	if sErr != nil {
		log.Errorf("RenewInstance err: %v", sErr)
		return xerrors.New(*sErr.Message)
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
		regions, sErr := aliyun.DescribeRegions(accessKeyID, accessKeySecret)
		if sErr != nil {
			log.Errorf("DescribePrice err:%v", sErr)
			return
		}
		for _, region := range regions.Body.Regions.Region {
			regionIDs = append(regionIDs, *region.RegionId)
		}
	}

	for _, regionID := range regionIDs {
		instanceType := &types.DescribeInstanceTypeReq{
			RegionId:     regionID,
			CpuCoreCount: 0,
			MemorySize:   0,
		}

		instances, err := m.DescribeInstanceType(ctx, instanceType)
		if err != nil {
			log.Errorf("DescribeInstanceType err:%v", err.Error())
			continue
		}
		for _, instance := range instances.InstanceTypes {
			ok, err := m.InstancesDefaultExists(instance.InstanceTypeId, regionID)
			if err != nil {
				log.Errorf("InstancesDefaultExists err:%v", err.Error())
				continue
			}
			if ok {
				continue
			}
			images, err := m.DescribeImages(ctx, regionID, instance.InstanceTypeId)
			if err != nil {
				log.Errorf("DescribePrice err:%v", err.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeId, regionID)
				continue
			}

			if len(images) == 0 {
				continue
			}

			disk := &types.AvailableResourceReq{
				InstanceType:        instance.InstanceTypeId,
				RegionId:            regionID,
				DestinationResource: "SystemDisk",
			}

			disks, err := m.DescribeAvailableResourceForDesk(ctx, disk)
			if err != nil {
				log.Errorf("DescribePrice err:%v", err.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeId, regionID)
				continue
			}

			if len(disks) == 0 {
				continue
			}

			priceReq := &types.DescribePriceReq{
				RegionId:                regionID,
				InstanceType:            instance.InstanceTypeId,
				PriceUnit:               "Month",
				ImageID:                 images[0].ImageId,
				InternetChargeType:      "PayByTraffic",
				SystemDiskCategory:      disks[0].Value,
				SystemDiskSize:          40,
				Period:                  1,
				Amount:                  1,
				InternetMaxBandwidthOut: 10,
			}

			price, sErr := aliyun.DescribePrice(accessKeyID, accessKeySecret, priceReq)
			if sErr != nil {
				log.Errorf("DescribePrice err:%v", sErr.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeId, regionID)
				continue
			}

			info := &types.DescribeInstanceTypeFromBase{
				RegionId:               regionID,
				InstanceTypeId:         instance.InstanceTypeId,
				MemorySize:             instance.MemorySize,
				CpuArchitecture:        instance.CpuArchitecture,
				InstanceCategory:       instance.InstanceCategory,
				CpuCoreCount:           instance.CpuCoreCount,
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
	return
}

// DescribeInstanceType fetches instance type information.
func (m *Manager) DescribeInstanceType(ctx context.Context, instanceType *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error) {
	accessKeyID := m.cfg.AliyunAccessKeyID
	accessKeySecret := m.cfg.AliyunAccessKeySecret
	rsp, sErr := aliyun.DescribeInstanceTypes(accessKeyID, accessKeySecret, instanceType)
	if sErr != nil {
		log.Errorf("DescribeInstanceTypes err: %v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}

	availableResource, sErr := aliyun.DescribeAvailableResource(accessKeyID, accessKeySecret, instanceType)
	if sErr != nil {
		log.Errorf("DescribeAvailableResource err: %v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
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
				InstanceTypeId:         *data.InstanceTypeId,
				MemorySize:             *data.MemorySize,
				CpuArchitecture:        *data.CpuArchitecture,
				InstanceTypeFamily:     *data.InstanceTypeFamily,
				CpuCoreCount:           *data.CpuCoreCount,
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

	rsp, sErr := aliyun.DescribeImages(regionID, accessKeyID, accessKeySecret, instanceType)
	if sErr != nil {
		log.Errorf("DescribeImages err: %v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}

	var rspDataList []*types.DescribeImageResponse
	for _, data := range rsp.Body.Images.Image {
		rspData := &types.DescribeImageResponse{
			ImageId:      *data.ImageId,
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
	rsp, sErr := aliyun.DescribeAvailableResourceForDesk(accessKeyID, accessKeySecret, desk)
	if sErr != nil {
		log.Errorf("DescribeImages err: %v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
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

	sErr := aliyun.ModifyInstanceAutoRenewAttribute(accessKeyID, accessKeySecret, renewReq)
	if sErr != nil {
		log.Errorf("ModifyInstanceAutoRenewAttribute err: %s", *sErr.Message)
		return &api.ErrWeb{Code: terrors.ThisInstanceNotSupportOperation.Int(), Message: *sErr.Message}
	}

	return nil
}
