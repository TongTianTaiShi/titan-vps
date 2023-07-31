package api

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
)

type Basis interface {
	Common

	DescribeRegions(ctx context.Context) ([]string, error)                                                                                      //perm:read
	DescribeInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error)                                   //perm:read
	DescribeImages(ctx context.Context, regionID, instanceType string) ([]string, error)                                                        //perm:read
	DescribePrice(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.DescribePriceResponse, error)   //perm:read
	CreateInstance(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.CreateInstanceResponse, error) //perm:read
	CreateKeyPair(ctx context.Context, regionID, KeyPairName string) (*types.CreateKeyPairResponse, error)                                      //perm:read
	AttachKeyPair(ctx context.Context, regionID, KeyPairName string, instanceIds []string) ([]*types.AttachKeyPairResponse, error)              //perm:read
	RebootInstance(ctx context.Context, regionID, instanceId string) (string, error)                                                            //perm:read
}
