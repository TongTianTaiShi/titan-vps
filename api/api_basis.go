package api

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
)

type Basis interface {
	Common

	Hello(ctx context.Context, id string) (*types.Hellos, error) //perm:read
}
