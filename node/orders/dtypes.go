package orders

import (
	"github.com/LMF709268224/titan-vps/api/types"
)

// OrderHash is an identifier for a asset.
type OrderHash string

func (c OrderHash) String() string {
	return string(c)
}

// OrderInfo represents asset pull information
type OrderInfo struct {
	State OrderState
	Hash  OrderHash
}

// ToOrderRecord converts AssetPullingInfo to types.AssetRecord
func (state *OrderInfo) ToOrderRecord() *types.OrderRecord {
	return &types.OrderRecord{
		Hash:  state.Hash.String(),
		State: state.State.String(),
	}
}

// orderInfoFrom converts types.AssetRecord to AssetPullingInfo
func orderInfoFrom(info *types.OrderRecord) *OrderInfo {
	cInfo := &OrderInfo{
		State: OrderState(info.State),
		Hash:  OrderHash(info.Hash),
	}

	return cInfo
}
