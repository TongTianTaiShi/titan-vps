package basis

import (
	"context"
	"math/big"
)

func (m *Basis) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	return m.TransactionMgr.GetBalance(address)
}

func (m *Basis) Recharge(ctx context.Context, address, rechargeAddr string) (string, error) {
	return m.ExchangeMgr.CreateRechargeOrder(address, rechargeAddr)
}
