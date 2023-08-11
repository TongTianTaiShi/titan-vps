package basis

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

var log = logging.Logger("basis")

// Basis represents a base service in a cloud computing system.
type Basis struct {
	fx.In
	*common.CommonAPI
	TransactionMgr *transaction.Manager
	*exchange.RechargeManager
	*exchange.WithdrawManager
	Notify *pubsub.PubSub
	*db.SQLDB
	OrderMgr *orders.Manager
	dtypes.GetBasisConfigFunc
	UserMgr *user.Manager
}

func (m *Basis) getAccessKeys() (string, string) {
	cfg, err := m.GetBasisConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", ""
	}

	return cfg.AliyunAccessKeyID, cfg.AliyunAccessKeySecret
}

func (m *Basis) DescribeRegions(ctx context.Context) ([]string, error) {
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

func (m *Basis) DescribeRecommendInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error) {
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

func (m *Basis) DescribeInstanceType(ctx context.Context, regionID, CpuArchitecture, InstanceCategory string, cores int32, memory float32) ([]types.DescribeInstanceTypeResponse, error) {
	k, s := m.getAccessKeys()
	rsp, err := aliyun.DescribeInstanceTypes(regionID, k, s, CpuArchitecture, InstanceCategory, cores, memory)
	if err != nil {
		log.Errorf("DescribeInstanceTypes err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}
	var rspDataList []types.DescribeInstanceTypeResponse
	for _, data := range rsp.Body.InstanceTypes.InstanceType {
		var rspData types.DescribeInstanceTypeResponse
		rspData.InstanceCategory = *data.InstanceCategory
		rspData.InstanceTypeId = *data.InstanceTypeId
		rspData.MemorySize = *data.MemorySize
		rspData.CpuArchitecture = *data.CpuArchitecture
		rspData.InstanceTypeFamily = *data.CpuArchitecture
		rspData.CpuCoreCount = *data.CpuCoreCount
		rspData.PhysicalProcessorModel = *data.PhysicalProcessorModel
		rspDataList = append(rspDataList, rspData)
	}
	return rspDataList, nil
}

func (m *Basis) DescribeImages(ctx context.Context, regionID, instanceType string) ([]string, error) {
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

func (m *Basis) DescribePrice(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.DescribePriceResponse, error) {
	k, s := m.getAccessKeys()

	price, err := aliyun.DescribePrice(regionID, k, s, instanceType, priceUnit, imageID, period)
	if err != nil {
		log.Errorf("DescribePrice err:", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return price, nil
}

func (m *Basis) CreateKeyPair(ctx context.Context, regionID, KeyPairName string) (*types.CreateKeyPairResponse, error) {
	k, s := m.getAccessKeys()

	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, KeyPairName)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return keyInfo, nil
}

func (m *Basis) AttachKeyPair(ctx context.Context, regionID, KeyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error) {
	k, s := m.getAccessKeys()
	AttachResult, err := aliyun.AttachKeyPair(regionID, k, s, KeyPairName, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	return AttachResult, nil
}

func (m *Basis) RebootInstance(ctx context.Context, regionID, instanceId string) (string, error) {
	k, s := m.getAccessKeys()
	RebootResult, err := aliyun.RebootInstance(regionID, k, s, instanceId)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return "", xerrors.New(*err.Data)
	}

	return RebootResult.Body.String(), nil
}

func (m *Basis) CreateInstance(ctx context.Context, vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	k, s := m.getAccessKeys()
	priceUnit := vpsInfo.PeriodUnit
	period := vpsInfo.Period
	regionID := vpsInfo.RegionId
	instanceType := vpsInfo.InstanceType
	imageID := vpsInfo.ImageId
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

	result, err := aliyun.CreateInstance(regionID, k, s, instanceType, imageID, securityGroupID, priceUnit, period, false)
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

func (m *Basis) MintToken(ctx context.Context, address string) (string, error) {
	cfg, err := m.GetBasisConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", err
	}

	valueStr := "9000000000000000000"

	client := filecoinbridge.NewGrpcClient(cfg.LotusHTTPSAddr, cfg.TitanContractorAddr)

	return client.Mint(cfg.PrivateKeyStr, address, valueStr)
}

func (m *Basis) SignCode(ctx context.Context, userId string) (string, error) {
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	verifyCode := "Vps(" + fmt.Sprintf("%06d", randNew.Intn(1000000)) + ")"
	m.UserMgr.User[userId] = &types.UserInfoTmp{}
	m.UserMgr.User[userId].UserLogin.SignCode = verifyCode
	m.UserMgr.User[userId].UserLogin.UserId = userId
	return verifyCode, nil
}

func (m *Basis) Login(ctx context.Context, user *types.UserReq) (*types.UserResponse, error) {
	userId := user.UserId
	code := m.UserMgr.User[userId].UserLogin.SignCode
	signature := user.Signature
	m.UserMgr.User[userId].UserLogin.SignCode = ""
	publicKey, err := VerifyMessage(code, signature)
	userId = strings.ToUpper(userId)
	publicKey = strings.ToUpper(publicKey)
	if publicKey != userId {
		return nil, err
	}
	p := types.JWTPayload{
		ID:    userId,
		Allow: []auth.Permission{api.RoleUser},
	}
	rsp := &types.UserResponse{}
	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return rsp, err
	}
	rsp.UserId = userId
	rsp.Token = string(tk)
	return rsp, nil
}

func (m *Basis) Logout(ctx context.Context, user *types.UserReq) error {
	nodeID := handler.GetID(ctx)
	log.Warnf("user id : %s", nodeID)
	// delete(m.UserMgr.User, user.UserId)
	return nil
}

var _ api.Basis = &Basis{}

// JwtPayload represents the payload of a JWT token.
type JwtPayload struct {
	Allow []auth.Permission
}

func VerifyMessage(message string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
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
		log.Errorf("Could not get a public get from the message signature")
	}
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}
