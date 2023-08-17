package api

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
)

type Basis interface {
	Common
	OrderAPI
	UserAPI
	AdminAPI

	DescribeRegions(ctx context.Context) ([]string, error)                                                                                                          //perm:default
	DescribeInstanceType(ctx context.Context, instanceTypeReq *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error)                          //perm:default
	DescribeRecommendInstanceType(ctx context.Context, instanceTypeReq *types.DescribeRecommendInstanceTypeReq) ([]*types.DescribeRecommendInstanceResponse, error) //perm:default
	DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error)                                                      //perm:default
	DescribePrice(ctx context.Context, describePriceReq *types.DescribePriceReq) (*types.DescribePriceResponse, error)                                              //perm:default
	CreateInstance(ctx context.Context, vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error)                                                    //perm:default
	CreateKeyPair(ctx context.Context, regionID, instanceID string) (*types.CreateKeyPairResponse, error)                                                           //perm:default
	AttachKeyPair(ctx context.Context, regionID, keyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error)                                  //perm:default
	RebootInstance(ctx context.Context, regionID, instanceID string) error                                                                                          //perm:default
	DescribeInstances(ctx context.Context, regionID, instanceId string) error                                                                                       //perm:default
}

// AdminAPI is an interface for admin
type AdminAPI interface {
	AddAdminUser(ctx context.Context, userID, nickName string) error                                //perm:admin
	GetAdminSignCode(ctx context.Context, userID string) (string, error)                            //perm:default
	LoginAdmin(ctx context.Context, user *types.UserReq) (*types.UserResponse, error)               //perm:default
	MintToken(ctx context.Context, address string) (string, error)                                  //perm:admin
	GetWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) //perm:default
	UpdateWithdrawalRecord(ctx context.Context, orderID, withdrawHash string) error                 //perm:admin
}

// OrderAPI is an interface for order
type OrderAPI interface {
	// order
	CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error)                           //perm:user
	GetOrderWaitingPayment(ctx context.Context, limit, offset int64) (*types.OrderRecordResponse, error) //perm:user
	GetOrderInfo(ctx context.Context, limit, offset int64) (*types.OrderRecordResponse, error)           //perm:user
	CancelOrder(ctx context.Context, orderID string) error                                               //perm:user
}

// UserAPI is an interface for user
type UserAPI interface {
	// user
	GetBalance(ctx context.Context) (string, error)                                                     //perm:user
	RebootInstance(ctx context.Context, regionID, instanceID string) error                              //perm:user
	GetSignCode(ctx context.Context, userID string) (string, error)                                     //perm:default
	Login(ctx context.Context, user *types.UserReq) (*types.UserResponse, error)                        //perm:default
	Logout(ctx context.Context, user *types.UserReq) error                                              //perm:user
	GetRechargeAddress(ctx context.Context) (string, error)                                             //perm:user
	Withdraw(ctx context.Context, withdrawAddr, value string) error                                     //perm:user
	GetUserRechargeRecords(ctx context.Context, limit, offset int64) (*types.RechargeResponse, error)   //perm:user
	GetUserWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) //perm:user
	GetUserInstanceRecords(ctx context.Context, limit, offset int64) (*types.MyInstanceResponse, error) //perm:user
	GetInstanceDetailsInfo(ctx context.Context, instanceID string) (*types.InstanceDetails, error)      //perm:user
}
