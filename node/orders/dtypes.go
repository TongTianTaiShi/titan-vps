package orders

import (
	"github.com/LMF709268224/titan-vps/api/types"
)

// OrderHash is an identifier for a order.
type OrderHash string

func (c OrderHash) String() string {
	return string(c)
}

type GoodsInfo struct {
	ID       string
	Password string
}

type PaymentInfo struct {
	ID    string
	From  string
	To    string
	Value int64
}

// OrderInfo represents order information
type OrderInfo struct {
	State         OrderState
	OrderID       OrderHash
	From          string
	To            string
	Value         int64
	DoneState     OrderDoneState
	CreatedHeight int64
	DoneHeight    int64
	VpsID         int64

	*PaymentInfo
	*GoodsInfo
}

// ToOrderRecord converts order info to types.orderRecord
func (state *OrderInfo) ToOrderRecord() *types.OrderRecord {
	return &types.OrderRecord{
		OrderID:       state.OrderID.String(),
		State:         state.State.Int(),
		From:          state.From,
		To:            state.To,
		Value:         state.Value,
		DoneState:     state.DoneState.Int(),
		DoneHeight:    state.DoneHeight,
		CreatedHeight: state.CreatedHeight,
		VpsID:         state.VpsID,
	}
}

// orderInfoFrom converts types.orderRecord to order info
func orderInfoFrom(info *types.OrderRecord) *OrderInfo {
	cInfo := &OrderInfo{
		State:         OrderState(info.State),
		OrderID:       OrderHash(info.OrderID),
		DoneState:     OrderDoneState(info.DoneState),
		DoneHeight:    info.DoneHeight,
		CreatedHeight: info.CreatedHeight,
		Value:         info.Value,
		From:          info.From,
		To:            info.To,
		VpsID:         info.VpsID,
	}

	return cInfo
}
