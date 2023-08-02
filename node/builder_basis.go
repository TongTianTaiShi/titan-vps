package node

import (
	"errors"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/filecoin"
	"github.com/LMF709268224/titan-vps/node/impl/basis"
	"github.com/LMF709268224/titan-vps/node/modules"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/orders"
	"github.com/LMF709268224/titan-vps/node/repo"
	"github.com/filecoin-project/pubsub"
	"go.uber.org/fx"

	"golang.org/x/xerrors"
)

func Basis(out *api.Basis) Option {
	return Options(
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the basis option must be set before Config option")),
		),

		func(s *Settings) error {
			s.nodeType = repo.Basis
			return nil
		},

		func(s *Settings) error {
			resAPI := &basis.Basis{}
			s.invokes[ExtractAPIKey] = fx.Populate(resAPI)
			*out = resAPI
			return nil
		},
	)
}

func ConfigBasis(c interface{}) Option {
	cfg, ok := c.(*config.BasisCfg)
	if !ok {
		return Error(xerrors.Errorf("invalid config from repo, got: %T", c))
	}

	return Options(
		Override(new(*config.BasisCfg), cfg),
		ConfigCommon(&cfg.Common),
		Override(new(dtypes.SetBasisConfigFunc), modules.NewSetBasisConfigFunc),
		Override(new(dtypes.GetBasisConfigFunc), modules.NewGetBasisConfigFunc),
		Override(new(*pubsub.PubSub), modules.NewPubSub),
		Override(new(dtypes.MetadataDS), modules.Datastore),
		Override(new(*db.SQLDB), modules.NewDB),
		Override(new(*filecoin.Manager), filecoin.NewManager),
		Override(new(*orders.Manager), modules.NewStorageManager),
	)
}
