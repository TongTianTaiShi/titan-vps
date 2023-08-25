package api

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
)

type Mall interface {
	Common
	OrderAPI
	UserAPI
	AdminAPI

	DescribeRegions(ctx context.Context) (map[string]string, error)                                                                                                 //perm:default
	UpdateInstanceDefaultInfo(ctx context.Context) error                                                                                                            //perm:default
	DescribeInstanceType(ctx context.Context, instanceTypeReq *types.DescribeInstanceTypeReq) (*types.DescribeInstanceTypeResponse, error)                          //perm:default
	DescribeRecommendInstanceType(ctx context.Context, instanceTypeReq *types.DescribeRecommendInstanceTypeReq) ([]*types.DescribeRecommendInstanceResponse, error) //perm:default
	DescribeImages(ctx context.Context, regionID, instanceType string) ([]*types.DescribeImageResponse, error)                                                      //perm:default
	DescribeAvailableResourceForDesk(ctx context.Context, desk *types.AvailableResourceReq) ([]*types.AvailableResourceResponse, error)                             //perm:default
	DescribePrice(ctx context.Context, describePriceReq *types.DescribePriceReq) (*types.DescribePriceResponse, error)                                              //perm:default
	CreateInstance(ctx context.Context, vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error)                                                    //perm:default
	CreateKeyPair(ctx context.Context, regionID, instanceID string) (*types.CreateKeyPairResponse, error)                                                           //perm:default
	AttachKeyPair(ctx context.Context, regionID, keyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error)                                  //perm:default
	RebootInstance(ctx context.Context, regionID, instanceID string) error                                                                                          //perm:default
	DescribeInstances(ctx context.Context, regionID, instanceId string) error                                                                                       //perm:default
	GetInstanceDefaultInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) (*types.InstanceTypeResponse, error)                                            //perm:default
	GetInstanceCpuInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) ([]*int32, error)                                                                   //perm:default
	GetInstanceMemoryInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) ([]*float32, error)                                                              //perm:default
	GetRenewInstance(ctx context.Context, renewReq types.SetRenewOrderReq) (string, error)                                                                          //perm:default
}

// AdminAPI is an interface for admin
type AdminAPI interface {
	AddAdminUser(ctx context.Context, userID, nickName string) error                                             //perm:admin
	GetAdminSignCode(ctx context.Context, userID string) (string, error)                                         //perm:default
	LoginAdmin(ctx context.Context, user *types.UserReq) (*types.LoginResponse, error)                           //perm:default
	GetWithdrawalRecords(ctx context.Context, req *types.GetWithdrawRequest) (*types.GetWithdrawResponse, error) //perm:default
	ApproveUserWithdrawal(ctx context.Context, orderID, withdrawHash string) error                               //perm:admin
	RejectUserWithdrawal(ctx context.Context, orderID string) error                                              //perm:admin
	GetRechargeAddresses(ctx context.Context, limit, page int64) (*types.GetRechargeAddressResponse, error)      //perm:admin
}

// OrderAPI is an interface for order
type OrderAPI interface {
	// order
	CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error)                             //perm:user
	RenewOrder(ctx context.Context, renewReq types.RenewOrderReq) (string, error)                          //perm:user
	RenewInstance(ctx context.Context, renewReq types.SetRenewOrderReq) error                              //perm:user
	GetUseWaitingPaymentOrders(ctx context.Context, limit, page int64) (*types.OrderRecordResponse, error) //perm:user
	GetUserOrderRecords(ctx context.Context, limit, page int64) (*types.OrderRecordResponse, error)        //perm:user
	CancelUserOrder(ctx context.Context, orderID string) error                                             //perm:user
	PaymentUserOrder(ctx context.Context, orderID string) error                                            //perm:user
}

// UserAPI is an interface for user
type UserAPI interface {
	// user
	GetBalance(ctx context.Context) (*types.UserInfo, error)                                             //perm:user
	RebootInstance(ctx context.Context, regionID, instanceID string) error                               //perm:user
	GetSignCode(ctx context.Context, userID string) (string, error)                                      //perm:default
	Login(ctx context.Context, user *types.UserReq) (*types.LoginResponse, error)                        //perm:default
	Logout(ctx context.Context, user *types.UserReq) error                                               //perm:user
	GetRechargeAddress(ctx context.Context) (string, error)                                              //perm:user
	Withdraw(ctx context.Context, withdrawAddr, value string) error                                      //perm:user
	GetUserRechargeRecords(ctx context.Context, limit, page int64) (*types.RechargeResponse, error)      //perm:user
	GetUserWithdrawalRecords(ctx context.Context, limit, page int64) (*types.GetWithdrawResponse, error) //perm:user
	GetUserInstanceRecords(ctx context.Context, limit, offset int64) (*types.MyInstanceResponse, error)  //perm:user
	GetInstanceDetailsInfo(ctx context.Context, instanceID string) (*types.InstanceDetails, error)       //perm:user
	UpdateInstanceName(ctx context.Context, instanceID, instanceName string) error                       //perm:user
}
