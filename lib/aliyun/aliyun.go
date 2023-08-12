package aliyun

import (
	"encoding/json"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/opentracing/opentracing-go/log"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	defaultRegionID = "cn-hangzhou"
)

// ClientInit /**
func newClient(regionID, keyID, keySecret string) (*ecs20140526.Client, *tea.SDKError) {
	configClient := &openapi.Config{
		AccessKeyId:     tea.String(keyID),
		AccessKeySecret: tea.String(keySecret),
	}

	configClient.RegionId = tea.String(regionID)
	client, err := ecs20140526.NewClient(configClient)
	if err != nil {
		errors := &tea.SDKError{}
		if _t, ok := err.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(err.Error())
		}
		return nil, errors
	}

	return client, nil
}

// CreateInstance crate an instance
func CreateInstance(keyID, keySecret string, instanceReq *types.CreateInstanceReq, dryRun bool) (*types.CreateInstanceResponse, *tea.SDKError) {
	var out *types.CreateInstanceResponse

	client, err := newClient(instanceReq.RegionId, keyID, keySecret)
	if err != nil {
		return out, err
	}

	createInstanceRequest := &ecs20140526.CreateInstanceRequest{
		RegionId:           tea.String(instanceReq.RegionId),
		InstanceType:       tea.String(instanceReq.InstanceType),
		DryRun:             tea.Bool(dryRun),
		ImageId:            tea.String(instanceReq.ImageId),
		SecurityGroupId:    tea.String(instanceReq.SecurityGroupId),
		InstanceChargeType: tea.String("PrePaid"),
		PeriodUnit:         tea.String(instanceReq.PeriodUnit),
		InternetChargeType: tea.String(instanceReq.InternetChargeType),
		Period:             tea.Int32(instanceReq.Period),
		//Password:                tea.String(password),
		InternetMaxBandwidthOut: tea.Int32(1),
		InternetMaxBandwidthIn:  tea.Int32(1),
		// TODO
		SystemDisk: &ecs20140526.CreateInstanceRequestSystemDisk{
			Size:     tea.Int32(instanceReq.SystemDiskSize),
			Category: tea.String(instanceReq.SystemDiskCategory),
		},
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, err := client.CreateInstanceWithOptions(createInstanceRequest, runtime)
		if err != nil {
			return err
		}

		out = &types.CreateInstanceResponse{
			InstanceID: *result.Body.InstanceId,
			OrderId:    *result.Body.OrderId,
			RequestId:  *result.Body.RequestId,
			TradePrice: *result.Body.TradePrice,
		}

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// RunInstances run aliyun instances
func RunInstances(regionID, keyID, keySecret, instanceType, imageID, password, securityGroupID, periodUnit string, period int32) (*types.CreateInstanceResponse, *tea.SDKError) {
	var out *types.CreateInstanceResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}

	createInstanceRequest := &ecs20140526.RunInstancesRequest{
		RegionId:           tea.String(regionID),
		InstanceType:       tea.String(instanceType),
		DryRun:             tea.Bool(true),
		ImageId:            tea.String(imageID),
		SecurityGroupId:    tea.String(securityGroupID),
		InstanceChargeType: tea.String("PrePaid"),
		PeriodUnit:         tea.String(periodUnit),
		Period:             tea.Int32(period),
		Password:           tea.String(password),
		// TODO
		InternetMaxBandwidthOut: tea.Int32(1),
		InternetMaxBandwidthIn:  tea.Int32(1),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, err := client.RunInstancesWithOptions(createInstanceRequest, runtime)
		if err != nil {
			return err
		}

		out = &types.CreateInstanceResponse{
			InstanceID: *result.Body.InstanceIdSets.InstanceIdSet[0],
			OrderId:    *result.Body.OrderId,
			RequestId:  *result.Body.RequestId,
			TradePrice: *result.Body.TradePrice,
		}

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// StartInstance start an instance
func StartInstance(regionID, keyID, keySecret, instanceID string) *tea.SDKError {
	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return err
	}

	startInstancesRequest := &ecs20140526.StartInstancesRequest{
		RegionId:   tea.String(regionID),
		InstanceId: tea.StringSlice([]string{instanceID}),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		_, err := client.StartInstancesWithOptions(startInstancesRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return errors
	}
	return nil
}

// DescribeSecurityGroups describe user security groups
func DescribeSecurityGroups(regionID, keyID, keySecret string) ([]string, *tea.SDKError) {
	var out []string

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}

	describeSecurityGroupsRequest := &ecs20140526.DescribeSecurityGroupsRequest{
		RegionId: tea.String(regionID),
		// NetworkType: tea.String("classic"),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		response, err := client.DescribeSecurityGroupsWithOptions(describeSecurityGroupsRequest, runtime)
		if err != nil {
			return err
		}

		grop := response.Body.SecurityGroups.SecurityGroup
		for _, g := range grop {
			out = append(out, *g.SecurityGroupId)
		}

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// DescribeInstanceAttribute describe attribute of instance
func DescribeInstanceAttribute(regionID, keyID, keySecret, instanceID string) (*ecs20140526.DescribeInstanceAttributeResponse, *tea.SDKError) {
	var out *ecs20140526.DescribeInstanceAttributeResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}

	describeInstanceAttributeRequest := &ecs20140526.DescribeInstanceAttributeRequest{
		InstanceId: tea.String(instanceID),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, err := client.DescribeInstanceAttributeWithOptions(describeInstanceAttributeRequest, runtime)
		if err != nil {
			return err
		}

		out = result

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// AllocatePublicIPAddress Allocate IP address
func AllocatePublicIPAddress(regionID, keyID, keySecret, instanceID string) (string, *tea.SDKError) {
	var out string

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}

	allocatePublicIPAddressRequest := &ecs20140526.AllocatePublicIpAddressRequest{
		InstanceId: tea.String(instanceID),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, err := client.AllocatePublicIpAddressWithOptions(allocatePublicIPAddressRequest, runtime)
		if err != nil {
			return err
		}

		out = *result.Body.IpAddress

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// DescribePrice describe instance price
func DescribePrice(keyID, keySecret string, priceReq *types.DescribePriceReq) (*types.DescribePriceResponse, *tea.SDKError) {
	var out *types.DescribePriceResponse

	client, err := newClient(priceReq.RegionId, keyID, keySecret)
	if err != nil {
		return out, err
	}

	describePriceRequest := &ecs20140526.DescribePriceRequest{
		RegionId:           tea.String(priceReq.RegionId),
		InstanceType:       tea.String(priceReq.InstanceType),
		PriceUnit:          tea.String(priceReq.PriceUnit),
		Period:             tea.Int32(priceReq.Period),
		ImageId:            tea.String(priceReq.ImageID),
		InternetChargeType: tea.String(priceReq.InternetChargeType),
		// todo 查询批量购买某种配置的云服务器ECS的价格
		Amount:                  tea.Int32(priceReq.Amount),
		InternetMaxBandwidthOut: tea.Int32(priceReq.InternetMaxBandwidthOut),
		// PayByBandwidth
		SystemDisk: &ecs20140526.DescribePriceRequestSystemDisk{
			Category: tea.String(priceReq.SystemDiskCategory),
			Size:     tea.Int32(priceReq.SystemDiskSize),
		},
		DataDisk: []*ecs20140526.DescribePriceRequestDataDisk{},
	}
	if len(priceReq.DataDisk) > 0 {
		for _, v := range priceReq.DataDisk {
			DataDiskInfo := &ecs20140526.DescribePriceRequestDataDisk{
				Category:         v.Category,
				PerformanceLevel: v.PerformanceLevel,
				Size:             v.Size,
			}
			describePriceRequest.DataDisk = append(describePriceRequest.DataDisk, DataDiskInfo)
		}
	}

	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, err := client.DescribePriceWithOptions(describePriceRequest, runtime)
		if err != nil {
			return err
		}
		price := result.Body.PriceInfo.Price
		out = &types.DescribePriceResponse{
			Currency:      *price.Currency,
			OriginalPrice: *price.OriginalPrice,
			TradePrice:    *price.TradePrice,
		}
		return nil
	}()
	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// AuthorizeSecurityGroup authorize security group
func AuthorizeSecurityGroup(regionID, keyID, keySecret, securityGroupID string) *tea.SDKError {
	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return err
	}

	authorizeSecurityGroupRequest := &ecs20140526.AuthorizeSecurityGroupRequest{
		RegionId:        tea.String(regionID),
		SecurityGroupId: tea.String(securityGroupID),
		Permissions: []*ecs20140526.AuthorizeSecurityGroupRequestPermissions{
			{
				// TODO
				IpProtocol:   tea.String("ALL"),
				SourceCidrIp: tea.String("0.0.0.0/0"),
				PortRange:    tea.String("-1/-1"),
			},
		},
	}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, err := client.AuthorizeSecurityGroupWithOptions(authorizeSecurityGroupRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()
	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return errors
	}
	return nil
}

// DescribeRegions describe regions
func DescribeRegions(keyID, keySecret string) (*ecs20140526.DescribeRegionsResponse, *tea.SDKError) {
	client, err := newClient(defaultRegionID, keyID, keySecret)
	if err != nil {
		return nil, err
	}

	var result *ecs20140526.DescribeRegionsResponse
	describeRegionsRequest := &ecs20140526.DescribeRegionsRequest{}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, _e = client.DescribeRegionsWithOptions(describeRegionsRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()
	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

// DescribeRecommendInstanceType Describe Instance Type
func DescribeRecommendInstanceType(regionID, keyID, keySecret string, cores int32, memory float32) (*ecs20140526.DescribeRecommendInstanceTypeResponse, *tea.SDKError) {
	var result *ecs20140526.DescribeRecommendInstanceTypeResponse
	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return result, err
	}

	describeRecommendInstanceTypeRequest := &ecs20140526.DescribeRecommendInstanceTypeRequest{
		NetworkType:        tea.String("vpc"),
		RegionId:           tea.String(regionID),
		Cores:              tea.Int32(cores),
		Memory:             tea.Float32(memory),
		InstanceChargeType: tea.String("PrePaid"),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.DescribeRecommendInstanceTypeWithOptions(describeRecommendInstanceTypeRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

func DescribeInstanceTypes(keyID, keySecret string, instanceType *types.DescribeInstanceTypeReq) (*ecs20140526.DescribeInstanceTypesResponse, *tea.SDKError) {
	var result *ecs20140526.DescribeInstanceTypesResponse
	client, err := newClient(instanceType.RegionId, keyID, keySecret)
	if err != nil {
		return result, err
	}
	describeInstanceTypesRequest := &ecs20140526.DescribeInstanceTypesRequest{}
	if instanceType.CpuArchitecture != "" {
		describeInstanceTypesRequest.CpuArchitecture = tea.String(instanceType.CpuArchitecture)
	}
	if instanceType.InstanceCategory != "" {
		describeInstanceTypesRequest.InstanceCategory = tea.String(instanceType.InstanceCategory)
	}
	if instanceType.CpuCoreCount != 0 {
		describeInstanceTypesRequest.MinimumCpuCoreCount = tea.Int32(instanceType.CpuCoreCount)
		describeInstanceTypesRequest.MaximumCpuCoreCount = tea.Int32(instanceType.CpuCoreCount)
	}
	if instanceType.MemorySize != 0 {
		describeInstanceTypesRequest.MinimumMemorySize = tea.Float32(instanceType.MemorySize)
		describeInstanceTypesRequest.MaximumMemorySize = tea.Float32(instanceType.MemorySize)
	}
	if instanceType.MaxResults != 0 {
		describeInstanceTypesRequest.MaxResults = tea.Int64(instanceType.MaxResults)
	}
	if instanceType.NextToken != "" {
		describeInstanceTypesRequest.NextToken = tea.String(instanceType.NextToken)
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.DescribeInstanceTypesWithOptions(describeInstanceTypesRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

// CreateSecurityGroup Create Security Group
func CreateSecurityGroup(regionID, keyID, keySecret string) (string, *tea.SDKError) {
	var out string

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}

	createSecurityGroupRequest := &ecs20140526.CreateSecurityGroupRequest{
		RegionId: tea.String(regionID),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, err := client.CreateSecurityGroupWithOptions(createSecurityGroupRequest, runtime)
		if err != nil {
			return err
		}

		out = *result.Body.SecurityGroupId
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// DescribeImages Describe Images
func DescribeImages(regionID, keyID, keySecret, instanceType string) (*ecs20140526.DescribeImagesResponse, *tea.SDKError) {
	var result *ecs20140526.DescribeImagesResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return result, err
	}

	createSecurityGroupRequest := &ecs20140526.DescribeImagesRequest{
		RegionId: tea.String(regionID),
	}
	if instanceType != "" {
		createSecurityGroupRequest.InstanceType = tea.String(instanceType)
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.DescribeImagesWithOptions(createSecurityGroupRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

// DescribeInstanceStatus query instance status
func DescribeInstanceStatus(regionID, keyID, keySecret string, InstanceId []string) (*ecs20140526.DescribeInstanceStatusResponse, *tea.SDKError) {
	var result *ecs20140526.DescribeInstanceStatusResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return result, err
	}

	createSecurityGroupRequest := &ecs20140526.DescribeInstanceStatusRequest{
		RegionId:   tea.String(regionID),
		InstanceId: tea.StringSlice(InstanceId),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.DescribeInstanceStatusWithOptions(createSecurityGroupRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

// DescribeInstances instance detail info
func DescribeInstances(regionID, keyID, keySecret string, InstanceIds []string) (*ecs20140526.DescribeInstancesResponse, *tea.SDKError) {
	var result *ecs20140526.DescribeInstancesResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return result, err
	}
	instanceIdsByte, e := json.Marshal(InstanceIds)
	if e != nil {
		log.Error(e)
	}
	instanceIdSting := string(instanceIdsByte)
	createSecurityGroupRequest := &ecs20140526.DescribeInstancesRequest{
		RegionId:    tea.String(regionID),
		InstanceIds: tea.String(instanceIdSting),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.DescribeInstancesWithOptions(createSecurityGroupRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

// DescribeAvailableResource Describe Resource
func DescribeAvailableResource(regionID, keyID, keySecret string, cores int32, memory float32) (*ecs20140526.DescribeAvailableResourceResponse, *tea.SDKError) {
	var result *ecs20140526.DescribeAvailableResourceResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return result, err
	}

	describeAvailableResourceRequest := &ecs20140526.DescribeAvailableResourceRequest{
		RegionId:            tea.String(regionID),
		DestinationResource: tea.String("InstanceType"),
		InstanceChargeType:  tea.String("PrePaid"),
		Cores:               tea.Int32(cores),
		Memory:              tea.Float32(memory),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.DescribeAvailableResourceWithOptions(describeAvailableResourceRequest, runtime)
		if _e != nil {
			return _e
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}

// CreateKeyPair Create key pair
func CreateKeyPair(regionID, keyID, keySecret, KeyPairName string) (*types.CreateKeyPairResponse, *tea.SDKError) {
	var out *types.CreateKeyPairResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}

	createKeyPairRequest := &ecs20140526.CreateKeyPairRequest{
		RegionId:    tea.String(regionID),
		KeyPairName: tea.String(KeyPairName),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e := client.CreateKeyPairWithOptions(createKeyPairRequest, runtime)
		if _e != nil {
			return _e
		}
		keyInfo := result.Body
		out = &types.CreateKeyPairResponse{
			KeyPairID:      *keyInfo.KeyPairId,
			KeyPairName:    *keyInfo.KeyPairName,
			PrivateKeyBody: *keyInfo.PrivateKeyBody,
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// AttachKeyPair Attach KeyPair
func AttachKeyPair(regionID, keyID, keySecret, KeyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, *tea.SDKError) {
	var out []*types.AttachKeyPairResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return out, err
	}
	instanceIdsByte, e := json.Marshal(instanceIds)
	if e != nil {
		log.Error(e)
	}
	instanceIdSting := string(instanceIdsByte)
	attachKeyPairRequest := &ecs20140526.AttachKeyPairRequest{
		RegionId:    tea.String(regionID),
		KeyPairName: tea.String(KeyPairName),
		// InstanceIds should be []string
		InstanceIds: tea.String(instanceIdSting),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e := client.AttachKeyPairWithOptions(attachKeyPairRequest, runtime)
		if _e != nil {
			return _e
		}
		for _, i := range result.Body.Results.Result {
			instanceInfo := &types.AttachKeyPairResponse{
				Code:       *i.Code,
				InstanceId: *i.InstanceId,
				Message:    *i.Message,
				Success:    *i.Success,
			}
			out = append(out, instanceInfo)
		}
		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// RebootInstance  Reboot Instance
func RebootInstance(regionID, keyID, keySecret, instanceId string) (*ecs20140526.RebootInstanceResponse, *tea.SDKError) {
	var result *ecs20140526.RebootInstanceResponse

	client, err := newClient(regionID, keyID, keySecret)
	if err != nil {
		return result, err
	}

	rebootInstanceRequest := &ecs20140526.RebootInstanceRequest{
		InstanceId: tea.String(instanceId),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _e = client.RebootInstanceWithOptions(rebootInstanceRequest, runtime)
		if _e != nil {
			return _e
		}

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return result, errors
	}
	return result, nil
}
