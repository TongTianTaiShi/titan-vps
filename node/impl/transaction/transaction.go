package transaction

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/node/common"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
)

var log = logging.Logger("transaction")

// Transaction represents a transaction service in a cloud computing system.
type Transaction struct {
	fx.In

	*common.CommonAPI
}

func (m *Transaction) Hello(ctx context.Context) error {
	// TODO implement me
	log.Infoln("implement me")
	return nil
}

var _ api.Transaction = &Transaction{}
