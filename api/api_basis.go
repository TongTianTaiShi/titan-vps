package api

import (
	"context"
	"math/big"

	"github.com/LMF709268224/titan-vps/api/types"
)

type Basis interface {
	Common
	OrderAPI
	UserAPI

	MintToken(ctx context.Context, address string) (string, error) //perm:read

	DescribeRegions(ctx context.Context) ([]string, error)                                                                                      //perm:read
	DescribeInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error)                                   //perm:read
	DescribeImages(ctx context.Context, regionID, instanceType string) ([]string, error)                                                        //perm:read
	DescribePrice(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.DescribePriceResponse, error)   //perm:read
	CreateInstance(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.CreateInstanceResponse, error) //perm:read
	CreateKeyPair(ctx context.Context, regionID, KeyPairName string) (*types.CreateKeyPairResponse, error)                                      //perm:read
	AttachKeyPair(ctx context.Context, regionID, KeyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error)              //perm:read
	RebootInstance(ctx context.Context, regionID, instanceId string) (string, error)                                                            //perm:read
}

// OrderAPI is an interface for order
type OrderAPI interface {
	CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error)           //perm:read
	PaymentCompleted(ctx context.Context, req types.PaymentCompletedReq) (string, error) //perm:read
	CancelOrder(ctx context.Context, orderID string) error                               //perm:read
}

// UserAPI is an interface for user
type UserAPI interface {
	GetBalance(ctx context.Context, address string) (*big.Int, error) //perm:read
}
