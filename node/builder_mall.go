package node

import (
	"errors"

	"github.com/LMF709268224/titan-vps/node/exchange"
	"github.com/LMF709268224/titan-vps/node/user"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/impl/mall"
	"github.com/LMF709268224/titan-vps/node/modules"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/orders"
	"github.com/LMF709268224/titan-vps/node/repo"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/filecoin-project/pubsub"
	"go.uber.org/fx"

	"golang.org/x/xerrors"
)

func Mall(out *api.Mall) Option {
	return Options(
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the mall option must be set before Config option")),
		),

		func(s *Settings) error {
			s.nodeType = repo.Mall
			return nil
		},

		func(s *Settings) error {
			resAPI := &mall.Mall{}
			s.invokes[ExtractAPIKey] = fx.Populate(resAPI)
			*out = resAPI
			return nil
		},
	)
}

func ConfigMall(c interface{}) Option {
	cfg, ok := c.(*config.MallCfg)
	if !ok {
		return Error(xerrors.Errorf("invalid config from repo, got: %T", c))
	}

	return Options(
		Override(new(*config.MallCfg), cfg),
		ConfigCommon(&cfg.Common),
		Override(new(dtypes.SetMallConfigFunc), modules.NewSetMallConfigFunc),
		Override(new(dtypes.GetMallConfigFunc), modules.NewGetMallConfigFunc),
		Override(new(*pubsub.PubSub), modules.NewPubSub),
		Override(new(dtypes.MetadataDS), modules.Datastore),
		Override(new(*db.SQLDB), modules.NewDB),
		Override(new(*transaction.Manager), transaction.NewManager),
		Override(new(*exchange.RechargeManager), exchange.NewRechargeManager),
		Override(new(*exchange.WithdrawManager), exchange.NewWithdrawManager),
		Override(new(*orders.Manager), modules.NewStorageManager),
		Override(new(*user.Manager), user.NewManager),
		Override(InitDataTables, db.InitTables),
	)
}
