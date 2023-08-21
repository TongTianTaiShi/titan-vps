package mall

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
)

func (m *Mall) CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error) {
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
		log.Errorf("DescribePrice:%v", err)
		return "", &api.ErrWeb{Code: terrors.DescribePriceError.Int(), Message: err.Error()}
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
		return "", err
	}
	TradePriceString := strconv.FormatFloat(float64(req.TradePrice), 'f', -1, 64)

	newBalanceString := strconv.FormatFloat(float64(priceInfo.USDPrice)*1000000000000000000, 'f', -1, 64)

	info := &types.OrderRecord{
		VpsID:      id,
		OrderID:    orderID,
		UserID:     userID,
		Value:      newBalanceString,
		TradePrice: TradePriceString,
	}

	err = m.OrderMgr.CreatedOrder(info)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

func (m *Mall) GetUseWaitingPaymentOrders(ctx context.Context, limit, offset int64) (*types.OrderRecordResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadOrderRecordByUserUndone(userID, limit, offset, m.OrderMgr.GetOrderTimeoutMinute())
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

func (m *Mall) GetUserOrderRecords(ctx context.Context, limit, offset int64) (*types.OrderRecordResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadOrderRecordsByUser(userID, limit, offset, m.OrderMgr.GetOrderTimeoutMinute())
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

func (m *Mall) CancelUserOrder(ctx context.Context, orderID string) error {
	userID := handler.GetID(ctx)
	return m.OrderMgr.CancelOrder(orderID, userID)
}

func (m *Mall) PaymentUserOrder(ctx context.Context, orderID string) error {
	userID := handler.GetID(ctx)
	return m.OrderMgr.PaymentCompleted(orderID, userID)
}
