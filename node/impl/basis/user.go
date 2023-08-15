package basis

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
)

func (m *Basis) GetBalance(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	return m.LoadUserBalance(userID)
}

func (m *Basis) GetRechargeAddress(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	addr, err := m.GetRechargeAddressOfUser(userID)
	if err != nil {
		return "", err
	}

	if addr == "" {
		return m.TransactionMgr.AllocateTronAddress(userID)
	}

	return addr, nil
}

func (m *Basis) Withdraw(ctx context.Context, withdrawAddr, value string) error {
	userID := handler.GetID(ctx)

	return m.WithdrawManager.CreateWithdrawOrder(userID, withdrawAddr, value)
}

func (m *Basis) GetRechargeRecord(ctx context.Context, limit, offset int64) (*types.RechargeResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadRechargeRecordsByUser(userID, limit, offset)
}

func (m *Basis) GetWithdrawRecord(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadWithdrawRecordsByUser(userID, limit, offset)
}
