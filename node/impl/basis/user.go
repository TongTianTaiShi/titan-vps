package basis

import (
	"context"
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
)

func (m *Basis) GetBalance(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	return m.LoadUserBalance(userID)
}

func (m *Basis) GetRechargeAddress(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)
	fmt.Println("----------------------------")
	fmt.Println(userID)
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

func (m *Basis) GetUserRechargeRecords(ctx context.Context, limit, offset int64) (*types.RechargeResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadRechargeRecordsByUser(userID, limit, offset)
}

func (m *Basis) GetUserWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadWithdrawRecordsByUser(userID, limit, offset)
}

func (m *Basis) GetUserInstanceRecords(ctx context.Context, limit, offset int64) (*types.MyInstanceResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadMyInstancesInfo(userID, limit, offset)
}
