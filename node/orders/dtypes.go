package orders

import (
	"github.com/LMF709268224/titan-vps/api/types"
)

// OrderHash is an identifier for a order.
type OrderHash string

func (c OrderHash) String() string {
	return string(c)
}

// GoodsInfo bug goods info
type GoodsInfo struct {
	ID       string
	Password string
}

// OrderInfo represents order information
type OrderInfo struct {
	State     OrderState
	OrderID   OrderHash
	User      string
	Value     string
	DoneState OrderDoneState
	OrderType int64
	VpsID     int64
	Msg       string
	EndTime   string

	*GoodsInfo
}

// ToOrderRecord converts order info to types.orderRecord
func (state *OrderInfo) ToOrderRecord() *types.OrderRecord {
	return &types.OrderRecord{
		OrderID:   state.OrderID.String(),
		State:     types.OrderState(state.State),
		UserID:    state.User,
		Value:     state.Value,
		DoneState: state.DoneState.Int(),
		VpsID:     state.VpsID,
		Msg:       state.Msg,
		EndTime:   state.EndTime,
		OrderType: types.OrderType(state.OrderType),
	}
}

// orderInfoFrom converts types.orderRecord to order info
func orderInfoFrom(info *types.OrderRecord) *OrderInfo {
	cInfo := &OrderInfo{
		State:     OrderState(info.State),
		OrderID:   OrderHash(info.OrderID),
		DoneState: OrderDoneState(info.DoneState),
		Value:     info.Value,
		VpsID:     info.VpsID,
		Msg:       info.Msg,
		User:      info.UserID,
		EndTime:   info.EndTime,
		OrderType: int64(info.OrderType),
	}
	return cInfo
}
