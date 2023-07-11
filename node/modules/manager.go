package modules

import (
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/repo"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("modules")

// NewSetTransactionConfigFunc creates a function to set the transaction config
func NewSetTransactionConfigFunc(r repo.LockedRepo) func(cfg config.TransactionCfg) error {
	return func(cfg config.TransactionCfg) (err error) {
		return r.SetConfig(func(raw interface{}) {
			_, ok := raw.(*config.TransactionCfg)
			if !ok {
				return
			}
		})
	}
}

// NewGetTransactionConfigFunc creates a function to get the transaction config
func NewGetTransactionConfigFunc(r repo.LockedRepo) func() (config.TransactionCfg, error) {
	return func() (out config.TransactionCfg, err error) {
		raw, err := r.Config()
		if err != nil {
			return
		}

		scfg, ok := raw.(*config.TransactionCfg)
		if !ok {
			return
		}

		out = *scfg
		return
	}
}
