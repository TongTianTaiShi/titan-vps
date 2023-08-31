package mall

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/LMF709268224/titan-vps/api/terrors"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/node/exchange"
	"github.com/LMF709268224/titan-vps/node/user"
	"github.com/LMF709268224/titan-vps/node/utils"
	"github.com/LMF709268224/titan-vps/node/vps"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
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
	VpsMgr  *vps.Manager
}

// getAliAccessKeys retrieves Aliyun access keys from the configuration.
func (m *Mall) getAliAccessKeys() (string, string) {
	cfg, err := m.GetMallConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", ""
	}

	return cfg.AliyunAccessKeyID, cfg.AliyunAccessKeySecret
}

// DescribeRegions retrieves and describes cloud regions.
func (m *Mall) DescribeRegions(ctx context.Context) (map[string]string, error) {
	rsp, sErr := aliyun.DescribeRegions(m.getAliAccessKeys())
	if sErr != nil {
		log.Errorf("DescribeRegions err: %v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}

	rpsData := make(map[string]string)
	// fmt.Printf("Response: %+v\n", response)
	for _, region := range rsp.Body.Regions.Region {
		// fmt.Printf("Region ID: %s\n", region.RegionId)
		// rpsData = append(rpsData, *region.RegionId)
		switch *region.RegionId {
		// Exclude specific regions
		case "ap-northeast-2", "ap-south-1", "eu-west-1", "ap-southeast-5", "ap-southeast-3",
			"s-east-1", "me-east-1", "us-east-1", "eu-central-1", "ap-northeast-1", "ap-southeast-2":
			continue
		}
		rpsData[*region.RegionId] = *region.LocalName
	}

	return rpsData, nil
}

// DescribeRecommendInstanceType retrieves recommended instance types.
func (m *Mall) DescribeRecommendInstanceType(ctx context.Context, instanceTypeReq *types.DescribeRecommendInstanceTypeReq) ([]*types.DescribeRecommendInstanceResponse, error) {
	startTime := time.Now()
	defer log.Debugf("DescribeRecommendInstanceType request time:%s", time.Since(startTime))

	accessKeyID, accessKeySecret := m.getAliAccessKeys()
	rsp, sErr := aliyun.DescribeRecommendInstanceType(accessKeyID, accessKeySecret, instanceTypeReq)
	if sErr != nil {
		log.Errorf("DescribeRecommendInstanceType err: %v", sErr)
		return nil, xerrors.New(*sErr.Message)
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

// DescribeInstanceType retrieves information about a specific instance type.
func (m *Mall) DescribeInstanceType(ctx context.Context, instanceType *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error) {
	startTime := time.Now()
	defer log.Debugf("DescribeInstanceType request time:%s", time.Since(startTime))

	return m.VpsMgr.DescribeInstanceType(ctx, instanceType)
}

// DescribeImages retrieves images for a specific region and instance type.
func (m *Mall) DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error) {
	startTime := time.Now()
	defer log.Debugf("DescribeImages request time:%s", time.Since(startTime))

	return m.VpsMgr.DescribeImages(ctx, regionID, instanceType)
}

// DescribeAvailableResourceForDesk retrieves available resources for a desk.
func (m *Mall) DescribeAvailableResourceForDesk(ctx context.Context, desk *types.AvailableResourceReq) ([]*types.AvailableResourceResponse, error) {
	startTime := time.Now()
	defer log.Debugf("DescribeAvailableResourceForDesk request time:%s", time.Since(startTime))

	return m.VpsMgr.DescribeAvailableResourceForDesk(ctx, desk)
}

// DescribePrice calculates the price for a specific configuration.
func (m *Mall) DescribePrice(ctx context.Context, priceReq *types.DescribePriceReq) (*types.DescribePriceResponse, error) {
	startTime := time.Now()
	defer log.Debugf("DescribePrice request time:%s", time.Since(startTime))
	log.Infof("DescribePrice :%v", priceReq)

	accessKeyID, accessKeySecret := m.getAliAccessKeys()

	price, sErr := aliyun.DescribePrice(accessKeyID, accessKeySecret, priceReq)
	if sErr != nil {
		log.Errorf("DescribePrice err:%v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}

	usdRate := utils.GetUSDRate()
	price.USDPrice = price.USDPrice / usdRate

	return price, nil
}

// CreateKeyPair creates a key pair for a region and instance.
func (m *Mall) CreateKeyPair(ctx context.Context, regionID, instanceID string) (*types.CreateKeyPairResponse, error) {
	startTime := time.Now()
	defer log.Debugf("CreateKeyPair request time:%s", time.Since(startTime))

	accessKeyID, accessKeySecret := m.getAliAccessKeys()
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))

	// TODO 密钥对有上限 不能无限创建
	keyPairNameNew := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	keyInfo, sErr := aliyun.CreateKeyPair(regionID, accessKeyID, accessKeySecret, keyPairNameNew)
	if sErr != nil {
		log.Errorf("CreateKeyPair err: %v", sErr)
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}

	var instanceIds []string
	instanceIds = append(instanceIds, instanceID)
	_, sErr = aliyun.AttachKeyPair(regionID, accessKeyID, accessKeySecret, keyPairNameNew, instanceIds)
	if sErr != nil {
		log.Errorf("AttachKeyPair err: %v", sErr)
	}

	go func() {
		time.Sleep(1 * time.Minute)
		sErr = aliyun.RebootInstance(regionID, accessKeyID, accessKeySecret, instanceID)
		if sErr != nil {
			log.Infoln("RebootInstance err:", sErr)
		}
	}()

	return keyInfo, nil
}

// RebootInstance reboots a specific instance.
func (m *Mall) RebootInstance(ctx context.Context, regionID, instanceID string) error {
	startTime := time.Now()
	defer log.Debugf("RebootInstance request time:%s", time.Since(startTime))

	accessKeyID, accessKeySecret := m.getAliAccessKeys()
	sErr := aliyun.RebootInstance(regionID, accessKeyID, accessKeySecret, instanceID)
	if sErr != nil {
		log.Errorf("AttachKeyPair err: %v", sErr)
		return &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}

	return nil
}

// DescribeInstances retrieves information about a specific instance.
func (m *Mall) DescribeInstances(ctx context.Context, regionID, instanceId string) error {
	startTime := time.Now()
	defer log.Debugf("DescribeInstances request time:%s", time.Since(startTime))

	accessKeyID, accessKeySecret := m.getAliAccessKeys()
	var instanceIds []string
	instanceIds = append(instanceIds, instanceId)
	_, sErr := aliyun.DescribeInstances(regionID, accessKeyID, accessKeySecret, instanceIds)
	if sErr != nil {
		log.Errorf("AttachKeyPair err: %v", sErr)
		return &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *sErr.Message}
	}
	return nil
}

// UpdateInstanceDefaultInfo updates default instance information for a region.
func (m *Mall) UpdateInstanceDefaultInfo(ctx context.Context, regionID string) error {
	go m.VpsMgr.UpdateInstanceDefaultInfo(regionID)
	return nil
}

// RenewInstance renews an instance with the provided renewal information.
func (m *Mall) RenewInstance(ctx context.Context, renewReq types.SetRenewOrderReq) error {
	startTime := time.Now()
	defer log.Debugf("RenewInstance request time:%s", time.Since(startTime))

	return m.VpsMgr.ModifyInstanceRenew(&renewReq)
}

// GetRenewInstance retrieves the renewal status for an instance.
func (m *Mall) GetRenewInstance(ctx context.Context, renewReq types.SetRenewOrderReq) (string, error) {
	accessKeyID, accessKeySecret := m.getAliAccessKeys()
	info, sErr := aliyun.DescribeInstanceAutoRenewAttribute(accessKeyID, accessKeySecret, &renewReq)
	if sErr != nil {
		log.Errorf("DescribeInstanceAutoRenewAttribute err: %v", sErr)
		return "", &api.ErrWeb{Code: terrors.ThisInstanceNotSupportOperation.Int(), Message: *sErr.Message}
	}
	out := *info.Body.InstanceRenewAttributes.InstanceRenewAttribute[0].RenewalStatus
	return out, nil
}

// verifyEthMessage verifies an Ethereum message signature.
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
	if err != nil {
		return "", err
	}

	if sigPublicKeyECDSA == nil {
		return "", xerrors.New("Could not get a public get from the message signature")
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}

var _ api.Mall = &Mall{}
