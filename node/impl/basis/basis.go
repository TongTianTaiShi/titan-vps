package basis

import (
	"context"
	"fmt"
	"github.com/LMF709268224/titan-vps/api"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/node/common"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/filecoin"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/orders"
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
	User        map[string]*types.UserInfoTmp
	FilecoinMgr *filecoin.Manager
	Notify      *pubsub.PubSub
	*db.SQLDB
	OrderMgr *orders.Manager
	dtypes.GetBasisConfigFunc
	Key *dtypes.APIAlg
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

func (m *Basis) DescribeInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error) {
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

func (m *Basis) CreateInstance(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.CreateInstanceResponse, error) {
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

	result, err := aliyun.CreateInstance(regionID, k, s, instanceType, imageID, securityGroupID, priceUnit, period, false)
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
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPairName := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, keyPairName)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
	} else {
		result.PrivateKey = keyInfo.PrivateKeyBody
	}
	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceId)
	_, err = aliyun.AttachKeyPair(regionID, k, s, keyPairName, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
	}
	go func() {
		time.Sleep(1 * time.Minute)

		err := aliyun.StartInstance(regionID, k, s, result.InstanceId)
		log.Infoln("StartInstance err:", err)
	}()

	return result, nil
}



func (m *Basis) SignCode(ctx context.Context, userId string) (string, error) {
	m.User[userId].UserLogin.SignCode = "abc123"
	return "abc123", nil
}

func (m *Basis) Login(ctx context.Context, user *types.UserReq) (*types.UserResponse, error) {
	userId := user.UserId
	code := m.User[userId].UserLogin.SignCode
	signature := user.Signature
	address := user.Address
	m.User[userId].UserLogin.SignCode = ""
	publicKey, err := VerifyMessage(code, signature)
	address = strings.ToUpper(address)
	publicKey = strings.ToUpper(publicKey)
	if publicKey != address {
		return nil, err
	}
	payload := &types.Token{
		UserId:     userId,
		Expiration: time.Now().Add(time.Hour),
	}
	var rsp *types.UserResponse
	tk, err := jwt.Sign(&payload, (*jwt.HMACSHA)(m.Key))
	if err != nil {
		return rsp, err
	}
	rsp.UserId = userId
	rsp.Token = string(tk)
	return rsp, nil
}

func (m *Basis) Logout(ctx context.Context, user *types.UserReq) error {
	delete(m.User, user.UserId)
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
