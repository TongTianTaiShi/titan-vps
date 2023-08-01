package modules

import (
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/repo"
	"github.com/filecoin-project/pubsub"
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

// NewSetBasisConfigFunc creates a function to set the basis config
func NewSetBasisConfigFunc(r repo.LockedRepo) func(cfg config.BasisCfg) error {
	return func(cfg config.BasisCfg) (err error) {
		return r.SetConfig(func(raw interface{}) {
			_, ok := raw.(*config.BasisCfg)
			if !ok {
				return
			}
		})
	}
}

// NewGetBasisConfigFunc creates a function to get the basis config
func NewGetBasisConfigFunc(r repo.LockedRepo) func() (config.BasisCfg, error) {
	return func() (out config.BasisCfg, err error) {
		raw, err := r.Config()
		if err != nil {
			return
		}

		scfg, ok := raw.(*config.BasisCfg)
		if !ok {
			return
		}

		out = *scfg
		return
	}
}

// NewPubSub returns a new pubsub instance with a buffer of 50
func NewPubSub() *pubsub.PubSub {
	return pubsub.New(50)
}

// NewDB returns an *sqlx.DB instance
func NewDB(cfg *config.BasisCfg) (*db.SQLDB, error) {
	return db.NewSQLDB(cfg.DatabaseAddress)
}
