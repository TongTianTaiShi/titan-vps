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

func (m *Manager) CreateAliYunInstance(vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	priceUnit := vpsInfo.PeriodUnit
	period := vpsInfo.Period
	regionID := vpsInfo.RegionId
	if priceUnit == "Year" {
		priceUnit = "Month"
		period = period * 12
	}

	var securityGroupID string

	group, err := aliyun.DescribeSecurityGroups(regionID, k, s)
	if err == nil && len(group) > 0 {
		securityGroupID = group[0]
	} else {
		securityGroupID, err = aliyun.CreateSecurityGroup(regionID, k, s)
		if err != nil {
			log.Errorf("CreateSecurityGroup err: %s", err.Error())
			return nil, xerrors.New(err.Error())
		}
	}
	log.Debugln("securityGroupID:", securityGroupID, " , DryRun:", vpsInfo.DryRun)
	result, err := aliyun.CreateInstance(k, s, vpsInfo, vpsInfo.DryRun)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(err.Error())
	}
	address, err := aliyun.AllocatePublicIPAddress(regionID, k, s, result.InstanceID)
	if err != nil {
		log.Errorf("AllocatePublicIpAddress err: %s", err.Error())
	} else {
		result.PublicIpAddress = address
	}

	err = aliyun.AuthorizeSecurityGroup(regionID, k, s, securityGroupID)
	if err != nil {
		log.Errorf("AuthorizeSecurityGroup err: %s", err.Error())
	}
	go func() {
		time.Sleep(1 * time.Minute)

		err = aliyun.StartInstance(regionID, k, s, result.InstanceID)
		if err != nil {
			log.Infoln("StartInstance err:", err)
		}
	}()
	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	instanceInfo, err := aliyun.DescribeInstances(regionID, k, s, instanceIds)
	if err != nil {
		log.Errorf("DescribeInstances err: %s", err.Error())
	}
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
		instanceDetailsInfo := &types.CreateInstanceReq{
			IpAddress:               *ip,
			InstanceId:              result.InstanceID,
			SecurityGroupId:         securityGroupId,
			OrderID:                 vpsInfo.OrderID,
			UserID:                  vpsInfo.UserID,
			OSType:                  *OSType,
			Cores:                   *Cores,
			Memory:                  float32(*Memory),
			InstanceName:            *InstanceName,
			ExpiredTime:             *ExpiredTime,
			InternetMaxBandwidthOut: *BandwidthOut,
		}
		errU := m.UpdateVpsInstance(instanceDetailsInfo)
		if errU != nil {
			log.Errorf("UpdateVpsInstance:%v", errU)
		}
	}
	return result, nil
}

func (m *Manager) RenewInstance(renewInstanceRequest *types.RenewInstanceRequest) error {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret

	err := aliyun.RenewInstance(k, s, renewInstanceRequest)
	if err != nil {
		log.Errorf("RenewInstance err: %s", err.Error())
		return xerrors.New(err.Error())
	}
	return nil
}

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

func (m *Manager) UpdateInstanceDefaultInfo(regionID string) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	var ctx context.Context

	var regionIDs []string
	if regionID != "" {
		regionIDs = append(regionIDs, regionID)
	} else {
		regions, err := aliyun.DescribeRegions(k, s)
		if err != nil {
			log.Errorf("DescribePrice err:%v", err.Error())
			return
		}
		for _, region := range regions.Body.Regions.Region {
			regionIDs = append(regionIDs, *region.RegionId)
		}
	}

	for _, regionId := range regionIDs {
		instanceType := &types.DescribeInstanceTypeReq{
			RegionId:     regionId,
			CpuCoreCount: 0,
			MemorySize:   0,
		}
		//if *region.RegionId != "eu-central-1" {
		//	continue
		//}
		instances, err := m.DescribeInstanceType(ctx, instanceType)
		if err != nil {
			log.Errorf("DescribeInstanceType err:%v", err.Error())
			continue
		}
		for _, instance := range instances.InstanceTypes {
			ok, err := m.InstancesDefaultExists(instance.InstanceTypeId, regionId)
			if err != nil {
				log.Errorf("InstancesDefaultExists err:%v", err.Error())
				continue
			}
			if ok {
				continue
			}
			images, err := m.DescribeImages(ctx, regionId, instance.InstanceTypeId)
			if err != nil {
				log.Errorf("DescribePrice err:%v", err.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeId, regionId)
				continue
			}
			disk := &types.AvailableResourceReq{
				InstanceType:        instance.InstanceTypeId,
				RegionId:            regionId,
				DestinationResource: "SystemDisk",
			}

			disks, err := m.DescribeAvailableResourceForDesk(ctx, disk)
			if err != nil {
				log.Errorf("DescribePrice err:%v", err.Error())
				_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeId, regionId)
				continue
			}
			if len(disks) > 0 {
				priceReq := &types.DescribePriceReq{
					RegionId:                regionId,
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
				price, err := aliyun.DescribePrice(k, s, priceReq)
				if err != nil {
					log.Errorf("DescribePrice err:%v", err.Error())
					_ = m.UpdateInstanceDefaultStatus(instance.InstanceTypeId, regionId)
					continue
				}
				info := &types.DescribeInstanceTypeFromBase{
					RegionId:               regionId,
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

	}
	return
}

func (m *Manager) DescribeInstanceType(ctx context.Context, instanceType *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	rsp, err := aliyun.DescribeInstanceTypes(k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeInstanceTypes err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}
	AvailableResource, err := aliyun.DescribeAvailableResource(k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeAvailableResource err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}
	rspDataList := &types.DescribeInstanceTypeResponse{
		NextToken: *rsp.Body.NextToken,
	}
	instanceTypes := make(map[string]string)
	if AvailableResource.Body.AvailableZones == nil {
		return nil, xerrors.New("parameter error")
	}
	AvailableZone := len(AvailableResource.Body.AvailableZones.AvailableZone)
	if AvailableZone < 0 {
		return rspDataList, nil
	}
	for _, data := range AvailableResource.Body.AvailableZones.AvailableZone {
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
				AvailableZone:          AvailableZone,
				PhysicalProcessorModel: *data.PhysicalProcessorModel,
				Status:                 v,
			}
			rspDataList.InstanceTypes = append(rspDataList.InstanceTypes, rspData)
		}
	}
	return rspDataList, nil
}

func (m *Manager) DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret

	rsp, err := aliyun.DescribeImages(regionID, k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
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

func (m *Manager) DescribeAvailableResourceForDesk(ctx context.Context, desk *types.AvailableResourceReq) ([]*types.AvailableResourceResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	rsp, err := aliyun.DescribeAvailableResourceForDesk(k, s, desk)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}
	Category := map[string]int{
		"cloud":            1,
		"cloud_essd":       1,
		"cloud_ssd":        1,
		"cloud_efficiency": 1,
		"ephemeral_ssd":    1,
	}
	var rspDataList []*types.AvailableResourceResponse
	if rsp.Body.AvailableZones == nil {
		log.Infoln(desk)
		return rspDataList, nil
	}
	if len(rsp.Body.AvailableZones.AvailableZone) > 0 {
		AvailableResources := rsp.Body.AvailableZones.AvailableZone[0].AvailableResources.AvailableResource
		if len(AvailableResources) > 0 {
			systemDesk := AvailableResources[0].SupportedResources.SupportedResource
			if len(systemDesk) > 0 {
				for _, data := range systemDesk {
					if *data.Status == "Available" {
						if _, ok := Category[*data.Value]; ok {
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

func (m *Manager) ModifyInstanceRenew(renewReq *types.SetRenewOrderReq) error {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret
	err := m.UpdateRenewInstanceStatus(renewReq)
	if err != nil {
		log.Errorf("UpdateRenewInstanceStatus:%v", err)
		return err
	}
	errSDK := aliyun.ModifyInstanceAutoRenewAttribute(k, s, renewReq)
	if err != nil {
		log.Errorf("ModifyInstanceAutoRenewAttribute err: %s", *errSDK.Message)
		return &api.ErrWeb{Code: terrors.ThisInstanceNotSupportOperation.Int(), Message: *errSDK.Message}
	}
	return nil
}
