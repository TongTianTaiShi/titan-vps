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

func (m *Mall) getAccessKeys() (string, string) {
	cfg, err := m.GetMallConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return "", ""
	}

	return cfg.AliyunAccessKeyID, cfg.AliyunAccessKeySecret
}

func (m *Mall) DescribeRegions(ctx context.Context) (map[string]string, error) {
	rsp, err := aliyun.DescribeRegions(m.getAccessKeys())
	if err != nil {
		log.Errorf("DescribeRegions err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	rpsData := make(map[string]string)
	// fmt.Printf("Response: %+v\n", response)
	for _, region := range rsp.Body.Regions.Region {
		// fmt.Printf("Region ID: %s\n", region.RegionId)
		// rpsData = append(rpsData, *region.RegionId)
		switch *region.RegionId {
		case "ap-northeast-2":
			continue
		case "ap-south-1":
			continue
		case "eu-west-1":
			continue
		case "ap-southeast-5":
			continue
		case "ap-southeast-3":
			continue
		case "s-east-1":
			continue
		case "me-east-1":
			continue
		case "us-east-1":
			continue
		case "eu-central-1":
			continue
		case "ap-northeast-1":
			continue
		case "ap-southeast-2":
			continue
		}
		rpsData[*region.RegionId] = *region.LocalName
	}

	return rpsData, nil
}

func (m *Mall) DescribeRecommendInstanceType(ctx context.Context, instanceTypeReq *types.DescribeRecommendInstanceTypeReq) ([]*types.DescribeRecommendInstanceResponse, error) {
	k, s := m.getAccessKeys()
	rsp, err := aliyun.DescribeRecommendInstanceType(k, s, instanceTypeReq)
	if err != nil {
		log.Errorf("DescribeRecommendInstanceType err: %s", err.Error())
		return nil, xerrors.New(err.Error())
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
	return m.VpsMgr.DescribeInstanceType(ctx, instanceType)
}

func (m *Mall) DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error) {
	return m.VpsMgr.DescribeImages(ctx, regionID, instanceType)
}

func (m *Mall) DescribeAvailableResourceForDesk(ctx context.Context, desk *types.AvailableResourceReq) ([]*types.AvailableResourceResponse, error) {
	return m.VpsMgr.DescribeAvailableResourceForDesk(ctx, desk)
}

func (m *Mall) DescribePrice(ctx context.Context, priceReq *types.DescribePriceReq) (*types.DescribePriceResponse, error) {
	k, s := m.getAccessKeys()

	price, err := aliyun.DescribePrice(k, s, priceReq)
	if err != nil {
		log.Errorf("DescribePrice err:%v", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	usdRate := utils.GetUSDRate()
	price.USDPrice = price.USDPrice / usdRate

	return price, nil
}

func (m *Mall) CreateKeyPair(ctx context.Context, regionID, instanceID string) (*types.CreateKeyPairResponse, error) {
	k, s := m.getAccessKeys()
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPairNameNew := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, keyPairNameNew)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
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

func (m *Mall) RebootInstance(ctx context.Context, regionID, instanceId string) error {
	k, s := m.getAccessKeys()
	err := aliyun.RebootInstance(regionID, k, s, instanceId)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}

	return nil
}

func (m *Mall) DescribeInstances(ctx context.Context, regionID, instanceId string) error {
	k, s := m.getAccessKeys()
	var instanceIds []string
	instanceIds = append(instanceIds, instanceId)
	_, err := aliyun.DescribeInstances(regionID, k, s, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
		return &api.ErrWeb{Code: terrors.AliApiGetFailed.Int(), Message: *err.Message}
	}
	return nil
}

func (m *Mall) UpdateInstanceDefaultInfo(ctx context.Context) error {
	go m.VpsMgr.UpdateInstanceDefaultInfo()
	return nil
}

func (m *Mall) RenewInstance(ctx context.Context, renewReq types.SetRenewOrderReq) error {
	return m.VpsMgr.ModifyInstanceRenew(&renewReq)
}

func (m *Mall) GetRenewInstance(ctx context.Context, renewReq types.SetRenewOrderReq) (string, error) {
	k, s := m.getAccessKeys()
	info, err := aliyun.DescribeInstanceAutoRenewAttribute(k, s, &renewReq)
	if err != nil {
		log.Errorf("DescribeInstanceAutoRenewAttribute err: %s", err.Error())
		return "", &api.ErrWeb{Code: terrors.ThisInstanceNotSupportOperation.Int(), Message: *err.Message}
	}
	out := *info.Body.InstanceRenewAttributes.InstanceRenewAttribute[0].RenewalStatus
	return out, nil
}

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
