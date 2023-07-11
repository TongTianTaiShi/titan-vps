package node

import (
	"errors"

	"github.com/LMF709268224/titan-vps/node/impl/transaction"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/repo"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

func Transaction(out *api.Transaction) Option {
	return Options(
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the transaction option must be set before Config option")),
		),

		func(s *Settings) error {
			s.nodeType = repo.Transaction
			return nil
		},

		func(s *Settings) error {
			resAPI := &transaction.Transaction{}
			s.invokes[ExtractAPIKey] = fx.Populate(resAPI)
			*out = resAPI
			return nil
		},
	)
}

func ConfigTransaction(c interface{}) Option {
	cfg, ok := c.(*config.TransactionCfg)
	if !ok {
		return Error(xerrors.Errorf("invalid config from repo, got: %T", c))
	}

	return Options(
		Override(new(*config.TransactionCfg), cfg),
		ConfigCommon(&cfg.Common),
	)
}
