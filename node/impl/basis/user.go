package basis

import (
	"context"
	"math/big"
)

func (m *Basis) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	return m.FilecoinMgr.GetBalance(address)
}
