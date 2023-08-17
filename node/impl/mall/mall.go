package mall

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/node/exchange"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/LMF709268224/titan-vps/node/user"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/node/common"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/orders"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/filecoin-project/pubsub"
	logging "github.com/ipfs/go-log/v2"

	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var log = logging.Logger("mall")

var USDRate float32

// Mall represents a base service in a cloud computing system.
type Mall struct {
	fx.In
	*common.CommonAPI
	TransactionMgr *transaction.Manager
	*exchange.RechargeManager
	*exchange.WithdrawManager
	Notify *pubsub.PubSub
	*db.SQLDB
	OrderMgr *orders.Manager
	dtypes.GetMallConfigFunc
	UserMgr *user.Manager
}

func (m *Mall) getAccessKeys() (string, string) {
	cfg, err := m.GetMallConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", ""
	}

	return cfg.AliyunAccessKeyID, cfg.AliyunAccessKeySecret
}

func (m *Mall) DescribeRegions(ctx context.Context) ([]string, error) {
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

func (m *Mall) DescribeRecommendInstanceType(ctx context.Context, instanceTypeReq *types.DescribeRecommendInstanceTypeReq) ([]*types.DescribeRecommendInstanceResponse, error) {
	k, s := m.getAccessKeys()
	rsp, err := aliyun.DescribeRecommendInstanceType(k, s, instanceTypeReq)
	if err != nil {
		log.Errorf("DescribeRecommendInstanceType err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	var rspDataList []*types.DescribeRecommendInstanceResponse
	for _, data := range rsp.Body.Data.RecommendInstanceType {
		rspData := &types.DescribeRecommendInstanceResponse{
			InstanceType:       *data.InstanceType.InstanceType,
			Memory:             *data.InstanceType.Memory,
			Cores:              *data.InstanceType.Cores,
			InstanceTypeFamily: *data.InstanceType.InstanceTypeFamily,
		}
		rspDataList = append(rspDataList, rspData)
	}

	return rspDataList, nil
}

func (m *Mall) DescribeInstanceType(ctx context.Context, instanceType *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error) {
	k, s := m.getAccessKeys()
	rsp, err := aliyun.DescribeInstanceTypes(k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeInstanceTypes err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	AvailableResource, err := aliyun.DescribeAvailableResource(k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeInstanceTypes err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	instanceTypes := make(map[string]int)
	for _, data := range AvailableResource.Body.AvailableZones.AvailableZone {
		availableTypes := data.AvailableResources.AvailableResource
		if len(availableTypes) > 0 {
			for _, instanceTypeResource := range availableTypes {
				Resources := instanceTypeResource.SupportedResources.SupportedResource
				if len(Resources) > 0 {
					for _, Resource := range Resources {
						if *Resource.Status == "Available" {
							instanceTypes[*Resource.Value] = 1
						}
					}
				}
			}
		}
	}
	rspDataList := &types.DescribeInstanceTypeResponse{
		NextToken: *rsp.Body.NextToken,
	}
	for _, data := range rsp.Body.InstanceTypes.InstanceType {
		if _, ok := instanceTypes[*data.InstanceTypeId]; !ok {
			continue
		}
		rspData := &types.DescribeInstanceType{
			InstanceCategory:       *data.InstanceCategory,
			InstanceTypeId:         *data.InstanceTypeId,
			MemorySize:             *data.MemorySize,
			CpuArchitecture:        *data.CpuArchitecture,
			InstanceTypeFamily:     *data.InstanceTypeFamily,
			CpuCoreCount:           *data.CpuCoreCount,
			PhysicalProcessorModel: *data.PhysicalProcessorModel,
		}
		rspDataList.InstanceTypes = append(rspDataList.InstanceTypes, rspData)
	}
	return rspDataList, nil
}

func (m *Mall) DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error) {
	k, s := m.getAccessKeys()

	rsp, err := aliyun.DescribeImages(regionID, k, s, instanceType)
	if err != nil {
		log.Errorf("DescribeImages err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
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

func (m *Mall) DescribePrice(ctx context.Context, priceReq *types.DescribePriceReq) (*types.DescribePriceResponse, error) {
	k, s := m.getAccessKeys()

	price, err := aliyun.DescribePrice(k, s, priceReq)
	if err != nil {
		log.Errorf("DescribePrice err:%v", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	UsdRate := aliyun.GetExchangeRate()
	if UsdRate > 0 {
		price.USDPrice = price.USDPrice / UsdRate
		USDRate = UsdRate
	} else {
		price.USDPrice = price.USDPrice / USDRate
	}

	return price, nil
}

func (m *Mall) CreateKeyPair(ctx context.Context, regionID, instanceID string) (*types.CreateKeyPairResponse, error) {
	k, s := m.getAccessKeys()
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPairNameNew := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, keyPairNameNew)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	var instanceIds []string
	instanceIds = append(instanceIds, instanceID)
	_, err = aliyun.AttachKeyPair(regionID, k, s, keyPairNameNew, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
	}
	go func() {
		time.Sleep(1 * time.Minute)
		err = aliyun.RebootInstance(regionID, k, s, instanceID)
		if err != nil {
			log.Infoln("RebootInstance err:", err)
		}
	}()
	return keyInfo, nil
}

func (m *Mall) AttachKeyPair(ctx context.Context, regionID, keyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error) {
	k, s := m.getAccessKeys()
	AttachResult, err := aliyun.AttachKeyPair(regionID, k, s, keyPairName, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return AttachResult, nil
}

func (m *Mall) RebootInstance(ctx context.Context, regionID, instanceId string) error {
	k, s := m.getAccessKeys()
	err := aliyun.RebootInstance(regionID, k, s, instanceId)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return xerrors.New(*err.Data)
	}

	return nil
}

func (m *Mall) DescribeInstances(ctx context.Context, regionID, instanceId string) error {
	k, s := m.getAccessKeys()
	var instanceIds []string
	instanceIds = append(instanceIds, instanceId)
	instanceInfo, err := aliyun.DescribeInstances(regionID, k, s, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return xerrors.New(*err.Data)
	}
	fmt.Println(instanceInfo.Body.Instances.Instance[0])
	fmt.Println(instanceInfo.Body.Instances.Instance[0].PublicIpAddress.IpAddress[0])
	return nil
}

func (m *Mall) CreateInstance(ctx context.Context, vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	k, s := m.getAccessKeys()
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
			return nil, xerrors.New(*err.Data)
		}
	}

	log.Debugf("securityGroupID : ", securityGroupID)
	result, err := aliyun.CreateInstance(k, s, vpsInfo, false)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
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
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPairName := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, keyPairName)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
	} else {
		result.PrivateKey = keyInfo.PrivateKeyBody
	}
	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	_, err = aliyun.AttachKeyPair(regionID, k, s, keyPairName, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
	}
	go func() {
		time.Sleep(1 * time.Minute)

		err := aliyun.StartInstance(regionID, k, s, result.InstanceID)
		log.Infoln("StartInstance err:", err)
	}()
	return result, nil
}

func (m *Mall) MintToken(ctx context.Context, address string) (string, error) {
	cfg, err := m.GetMallConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", err
	}

	valueStr := "9000000000000000000"

	client := filecoinbridge.NewGrpcClient(cfg.LotusHTTPSAddr, cfg.TitanContractorAddr)

	return client.Mint(cfg.PrivateKeyStr, address, valueStr)
}

func (m *Mall) GetSignCode(ctx context.Context, userID string) (string, error) {
	return m.UserMgr.SetSignCode(userID)
}

func (m *Mall) Login(ctx context.Context, user *types.UserReq) (*types.UserResponse, error) {
	userID := user.UserId
	code, err := m.UserMgr.GetSignCode(userID)
	if err != nil {
		return nil, err
	}
	signature := user.Signature
	address, err := verifyEthMessage(code, signature)
	if err != nil {
		return nil, err
	}

	p := types.JWTPayload{
		ID:        address,
		LoginType: int64(user.Type),
		Allow:     []auth.Permission{api.RoleUser},
	}
	rsp := &types.UserResponse{}
	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return rsp, err
	}
	rsp.UserId = address
	rsp.Token = string(tk)

	err = m.initUser(address)
	if err != nil {
		return nil, xerrors.Errorf("initUser err:%s", err.Error())
	}

	return rsp, nil
}

func (m *Mall) initUser(userID string) error {
	exist, err := m.UserExists(userID)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	return m.SaveUserInfo(&types.UserInfo{UserID: userID, Balance: "0"})
}

func (m *Mall) Logout(ctx context.Context, user *types.UserReq) error {
	userID := handler.GetID(ctx)
	log.Warnf("user id : %s", userID)
	// delete(m.UserMgr.User, user.UserId)
	return nil
}

var _ api.Mall = &Mall{}

func verifyEthMessage(code string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(code)) + code)
	hash := crypto.Keccak256Hash(hashedMessage)
	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)
	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}
	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		return "", xerrors.New("Could not get a public get from the message signature")
	}
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}
