package basis

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/node/handler"
)

func (m *Basis) GetBalance(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	cfg, err := m.GetBasisConfigFunc()
	if err != nil {
		return "", err
	}

	client := filecoinbridge.NewGrpcClient(cfg.LotusHTTPSAddr, cfg.TitanContractorAddr)

	value, err := client.GetBalance(userID)
	if err != nil {
		return "", err
	}

	return value.String(), nil
}

func (m *Basis) Recharge(ctx context.Context, rechargeAddr string) (string, error) {
	userID := handler.GetID(ctx)

	return m.RechargeManager.CreateRechargeOrder(userID, rechargeAddr)
}

func (m *Basis) CancelRecharge(ctx context.Context, orderID string) error {
	return m.RechargeManager.CancelRechargeOrder(orderID)
}

func (m *Basis) Withdraw(ctx context.Context, withdrawAddr string) (string, error) {
	userID := handler.GetID(ctx)

	return m.WithdrawManager.CreateWithdrawOrder(userID, withdrawAddr)
}

func (m *Basis) CancelWithdraw(ctx context.Context, orderID string) error {
	return m.WithdrawManager.CancelWithdrawOrder(orderID)
}

func (m *Basis) GetRechargeRecord(ctx context.Context, limit, offset int64) ([]*types.RechargeRecord, error) {
	userID := handler.GetID(ctx)

	return m.LoadRechargeRecordsByUser(userID, limit, offset)
}

func (m *Basis) GetWithdrawRecord(ctx context.Context, limit, offset int64) ([]*types.WithdrawRecord, error) {
	userID := handler.GetID(ctx)

	return m.LoadWithdrawRecordsByUser(userID, limit, offset)
}
