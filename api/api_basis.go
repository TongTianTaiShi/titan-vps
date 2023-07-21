package api

import (
	"context"

	"github.com/LMF709268224/titan-vps/api/types"
)

type Basis interface {
	Common

	DescribeRegions(ctx context.Context) ([]string, error)                                                                                                //perm:read
	DescribeInstanceType(ctx context.Context, regionID string, cores int32, memory float32) ([]string, error)                                             //perm:read
	DescribeImages(ctx context.Context, regionID, instanceType string) ([]string, error)                                                                  //perm:read
	DescribePrice(ctx context.Context, regionID, instanceType, priceUnit, imageID string, period int32) (*types.DescribePriceResponse, error)             //perm:read
	CreateInstance(ctx context.Context, regionID, instanceType, priceUnit, imageID, password string, period int32) (*types.CreateInstanceResponse, error) //perm:read
}
