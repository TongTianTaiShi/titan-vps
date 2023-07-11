package basis

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/common"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
)

var log = logging.Logger("basis")

// Basis represents a base service in a cloud computing system.
type Basis struct {
	fx.In

	*common.CommonAPI
}

func (p *Basis) Hello(ctx context.Context, id string) (*types.Hellos, error) {
	// TODO implement me
	log.Infoln("implement me")
	return &types.Hellos{Msg: "hello"}, nil
}

var _ api.Basis = &Basis{}
