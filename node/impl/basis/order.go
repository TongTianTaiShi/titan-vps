package basis

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
)

func (m *Basis) CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error) {
	// instanceInfo := req.Vps

	info := &types.OrderRecord{
		VpsID: "abc",
		From:  req.User,
		Value: 10,
	}

	err := m.OrderMgr.CreatedOrder(info)
	if err != nil {
		return "", err
	}

	return info.To, nil
}

func (m *Basis) PaymentCompleted(ctx context.Context, req types.PaymentCompletedReq) (string, error) {
	return "", m.FilecoinMgr.CheckMessage(req.TransactionID)
}

func (m *Basis) CancelOrder(ctx context.Context, orderID string) error {
	return m.OrderMgr.CancelOrder(orderID)
}
