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
	State         OrderState
	OrderID       OrderHash
	From          string
	To            string
	Value         int64
	DoneState     int64
	CreatedHeight int64
	DoneHeight    int64
	VpsID         string
}

// ToOrderRecord converts AssetPullingInfo to types.AssetRecord
func (state *OrderInfo) ToOrderRecord() *types.OrderRecord {
	return &types.OrderRecord{
		OrderID:       state.OrderID.String(),
		State:         int64(state.State),
		From:          state.From,
		To:            state.To,
		Value:         state.Value,
		DoneState:     state.DoneState,
		DoneHeight:    state.DoneHeight,
		CreatedHeight: state.CreatedHeight,
		VpsID:         state.VpsID,
	}
}

// orderInfoFrom converts types.AssetRecord to AssetPullingInfo
func orderInfoFrom(info *types.OrderRecord) *OrderInfo {
	cInfo := &OrderInfo{
		State:         OrderState(info.State),
		OrderID:       OrderHash(info.OrderID),
		DoneState:     info.DoneState,
		DoneHeight:    info.DoneHeight,
		CreatedHeight: info.CreatedHeight,
		Value:         info.Value,
		From:          info.From,
		To:            info.To,
		VpsID:         info.VpsID,
	}

	return cInfo
}
