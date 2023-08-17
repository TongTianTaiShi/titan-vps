package basis

import (
	"context"
	"github.com/LMF709268224/titan-vps/node/utils"
	"github.com/google/uuid"
	"strconv"
	"strings"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

func (m *Basis) CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error) {
	userID := handler.GetID(ctx)
	priceReq := &types.DescribePriceReq{
		RegionId:                req.RegionId,
		InstanceType:            req.InstanceType,
		PriceUnit:               req.PeriodUnit,
		Period:                  req.Period,
		Amount:                  req.Amount,
		InternetChargeType:      req.InternetChargeType,
		ImageID:                 req.ImageId,
		InternetMaxBandwidthOut: req.InternetMaxBandwidthOut,
		SystemDiskCategory:      req.SystemDiskCategory,
		SystemDiskSize:          req.SystemDiskSize,
	}
	priceInfo, err := m.DescribePrice(ctx, priceReq)
	if err != nil {
		return "", err
	}

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)
	req.OrderID = orderID
	req.TradePrice = priceInfo.USDPrice / float32(req.Amount)
	//for i := int32(0); i < req.Amount; i++ {
	//	id, err = m.SaveVpsInstance(&req)
	//	if err != nil {
	//		log.Errorf("SaveVpsInstance:%v", err)
	//	}
	//}
	id, err := m.SaveVpsInstance(&req)
	if err != nil {
		log.Errorf("SaveVpsInstance:%v", err)
	}
	info := &types.OrderRecord{
		VpsID:      id,
		OrderID:    orderID,
		UserID:     userID,
		Value:      "10000000000",
		TradePrice: priceInfo.TradePrice,
	}
	oldBalance, err := m.LoadUserBalance(userID)
	if err != nil {
		log.Errorf("LoadUserBalance:%v", err)
	}
	newBalanceString := strconv.FormatFloat(float64(priceInfo.USDPrice)*1000000000000000000, 'f', -1, 64)
	newBalanceString, ok := utils.BigIntReduce(oldBalance, newBalanceString)
	if ok {
		err = m.UpdateUserBalance(userID, newBalanceString, oldBalance)
		if err != nil {
			log.Errorf("UpdateUserBalance:%v", err)
			return "", err
		}
	}
	err = m.OrderMgr.CreatedOrder(info)
	if err != nil {
		return "", err
	}

	return info.To, nil
}

func (m *Basis) GetOrderWaitingPayment(ctx context.Context, limit, offset int64) (*types.OrderRecordResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadOrderRecordByUserUndone(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (m *Basis) GetOrderInfo(ctx context.Context, limit, offset int64) (*types.OrderRecordResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadOrderRecordByUserAll(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (m *Basis) PaymentCompleted(ctx context.Context, req types.PaymentCompletedReq) (string, error) {
	record, err := m.LoadOrderRecord(req.OrderID)
	if err != nil {
		return "", err
	}

	if record.State != types.WaitingPayment {
		return "", xerrors.Errorf("Invalid order status %d", record.State)
	}

	cfg, err := m.GetBasisConfigFunc()
	if err != nil {
		return "", err
	}

	tx := req.TransactionID
	log.Debugf("tx:%s \n", tx)
	var cid cid.Cid
	err = filecoinbridge.EthGetMessageCidByTransactionHash(&cid, tx, cfg.LotusHTTPSAddr)
	if err != nil {
		return "", err
	}

	log.Debugf("cid:%s \n", cid.String())

	var msg filecoinbridge.Message
	err = filecoinbridge.ChainGetMessage(&msg, cid, cfg.LotusHTTPSAddr)
	if err != nil {
		return "", err
	}

	var info filecoinbridge.Lookup
	err = filecoinbridge.StateSearchMsg(&info, cid, cfg.LotusHTTPSAddr)
	if err != nil {
		return "", err
	}

	log.Debugf("Height:%d,ExitCode:%d,GasUsed:%d \n", info.Height, info.Receipt.ExitCode, info.Receipt.GasUsed)

	if info.Receipt.ExitCode == 0 {
		m.Notify.Pub(&types.FvmTransferWatch{
			TxHash: tx,
			From:   msg.From.String(),
			To:     msg.To.String(),
			Value:  msg.Value.String(),
		}, types.EventFvmTransferWatch.String())
	}

	return "", nil
}

func (m *Basis) CancelOrder(ctx context.Context, orderID string) error {
	return m.OrderMgr.CancelOrder(orderID)
}
