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

	MintToken(ctx context.Context, address string) (string, error) //perm:default

	DescribeRegions(ctx context.Context) ([]string, error)                                                                                    //perm:default
	DescribeInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error)                                 //perm:default
	DescribeImages(ctx context.Context, regionID, instanceType string) ([]string, error)                                                      //perm:default
	DescribePrice(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.DescribePriceResponse, error) //perm:default
	CreateInstance(ctx context.Context, vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error)                              //perm:default
	CreateKeyPair(ctx context.Context, regionID, KeyPairName string) (*types.CreateKeyPairResponse, error)                                    //perm:default
	AttachKeyPair(ctx context.Context, regionID, KeyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error)            //perm:default
	RebootInstance(ctx context.Context, regionID, instanceID string) (string, error)                                                          //perm:default
}

// OrderAPI is an interface for order
type OrderAPI interface {
	// order
	CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error)           //perm:default
	PaymentCompleted(ctx context.Context, req types.PaymentCompletedReq) (string, error) //perm:default
	CancelOrder(ctx context.Context, orderID string) error                               //perm:default
}

// UserAPI is an interface for user
type UserAPI interface {
	// user
	GetBalance(ctx context.Context, address string) (*big.Int, error)                //perm:default
	RebootInstance(ctx context.Context, regionID, instanceID string) (string, error) //perm:default
	SignCode(ctx context.Context, userID string) (string, error)                     //perm:default
	Login(ctx context.Context, user *types.UserReq) (*types.UserResponse, error)     //perm:default
	Logout(ctx context.Context, user *types.UserReq) error                           //perm:default
	Recharge(ctx context.Context, address, rechargeAddr string) (string, error)      //perm:default
	CancelRecharge(ctx context.Context, orderID string) error                        //perm:default
	Withdraw(ctx context.Context, address, withdrawAddr string) (string, error)      //perm:default
	CancelWithdraw(ctx context.Context, orderID string) error                        //perm:default
}
