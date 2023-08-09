package basis

import (
	"context"
	"math/big"

	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
)

func (m *Basis) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	cfg, err := m.GetBasisConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return big.NewInt(0), err
	}

	client := filecoinbridge.NewGrpcClient(cfg.LotusHTTPSAddr, cfg.TitanContractorAddr)

	return client.GetBalance(address)
}

func (m *Basis) Recharge(ctx context.Context, address, rechargeAddr string) (string, error) {
	return m.ExchangeMgr.CreateRechargeOrder(address, rechargeAddr)
}

func (m *Basis) CancelRecharge(ctx context.Context, orderID string) error {
	return m.ExchangeMgr.CancelRechargeOrder(orderID)
}
