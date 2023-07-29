package basis

import (
	"context"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/node/common"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var log = logging.Logger("basis")

// Manager represents a base service in a cloud computing system.
type Manager struct {
	fx.In

	*common.CommonAPI

	GetConfigFunc dtypes.GetBasisConfigFunc
}

func (m *Manager) getAccessKeys() (string, string) {
	cfg, err := m.GetConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", ""
	}

	return cfg.AliyunAccessKeyID, cfg.AliyunAccessKeySecret
}

func (m *Manager) DescribeRegions(ctx context.Context) ([]string, error) {
	rsp, err := aliyun.DescribeRegions(m.getAccessKeys())
	if err != nil {
		log.Errorf("DescribeRegions err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	var rpsData []string
	// fmt.Printf("Response: %+v\n", response)
	for _, region := range rsp.Body.Regions.Region {
		// fmt.Printf("Region ID: %s\n", region.RegionId)
		rpsData = append(rpsData, *region.RegionId)
	}

	return rpsData, nil
}

func (m *Manager) DescribeInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error) {
	k, s := m.getAccessKeys()

	rsp, err := aliyun.DescribeRecommendInstanceType(regionID, k, s, cores, memory)
	if err != nil {
		log.Errorf("DescribeRecommendInstanceType err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	resources := make(map[string]string)
	for _, data := range rsp.Body.Data.RecommendInstanceType {
		instanceType := data.InstanceType.InstanceType
		if *instanceType == "" {
			continue
		}
		resources[*instanceType] = *instanceType
	}

	var rpsData []string
	for value := range resources {
		rpsData = append(rpsData, value)
	}

	return rpsData, nil
}

func (m *Manager) DescribeImages(ctx context.Context, regionID, instanceType string) ([]string, error) {
	k, s := m.getAccessKeys()

	rsp, err := aliyun.DescribeImages(regionID, k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	var rpsData []string
	for _, data := range rsp.Body.Images.Image {
		instanceType := data.ImageId
		if *instanceType == "" {
			continue
		}
		rpsData = append(rpsData, *instanceType)
	}

	return rpsData, nil
}

func (m *Manager) DescribePrice(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.DescribePriceResponse, error) {
	k, s := m.getAccessKeys()

	price, err := aliyun.DescribePrice(regionID, k, s, instanceType, priceUnit, imageID, period)
	if err != nil {
		log.Errorf("DescribePrice err:", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return price, nil
}

func (m *Manager) CreateKeyPair(ctx context.Context, regionID, KeyPairName string) (*types.CreateKeyPairResponse, error) {
	k, s := m.getAccessKeys()

	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, KeyPairName)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return keyInfo, nil
}

func (m *Manager) AttachKeyPair(ctx context.Context, regionID, KeyPairName, instanceIds string) ([]*types.AttachKeyPairResponse, error) {
	k, s := m.getAccessKeys()
	AttachResult, err := aliyun.AttachKeyPair(regionID, k, s, KeyPairName, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return AttachResult, nil
}

func (m *Manager) RebootInstance(ctx context.Context, regionID, instanceId string) (string, error) {
	k, s := m.getAccessKeys()
	RebootResult, err := aliyun.RebootInstance(regionID, k, s, instanceId)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return "", xerrors.New(*err.Data)
	}

	return RebootResult.Body.String(), nil
}

func (m *Manager) CreateInstance(ctx context.Context, regionID, instanceType, priceUnit, imageID, password string, period int32) (*types.CreateInstanceResponse, error) {
	k, s := m.getAccessKeys()

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
			return nil, xerrors.New(*err.Data)
		}
	}

	log.Debugf("securityGroupID : ", securityGroupID)

	result, err := aliyun.CreateInstance(regionID, k, s, instanceType, imageID, password, securityGroupID, priceUnit, period, false)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	address, err := aliyun.AllocatePublicIPAddress(regionID, k, s, result.InstanceId)
	if err != nil {
		log.Errorf("AllocatePublicIpAddress err: %s", err.Error())
	} else {
		result.PublicIpAddress = address
	}

	err = aliyun.AuthorizeSecurityGroup(regionID, k, s, securityGroupID)
	if err != nil {
		log.Errorf("AuthorizeSecurityGroup err: %s", err.Error())
	}

	// 一分钟后调用
	go func() {
		time.Sleep(1 * time.Minute)

		err := aliyun.StartInstance(regionID, k, s, result.InstanceId)
		log.Infoln("StartInstance err:", err)
	}()

	return result, nil
}

var _ api.Basis = &Manager{}
