package orders

// OrderState represents the different states of an order during processing.
type OrderState int64

// Constants defining various states of the order during processing.
const (
	// OrderStateCreated represents the state when the order is created
	OrderStateCreated OrderState = iota
	// OrderStateWaitingPayment represents the state when the order is waiting for user payment.
	OrderStateWaitingPayment
	// OrderStateBuyGoods represents the state when the order is in the process of buying goods.
	OrderStateBuyGoods
	// OrderStateDone represents the state when the order is completed.
	OrderStateDone
)

// String returns the string representation of the order state.
func (s OrderState) String() string {
	switch s {
	case 0:
		return "Created"
	case 1:
		return "WaitingPayment"
	case 2:
		return "BuyGoods"
	case 3:
		return "Done"
	}

	return "Not found"
}

// Int returns the int representation of the order state.
func (s OrderState) Int() int64 {
	return int64(s)
}

var (
	// ActiveStates contains a list of order states that represent active orders.
	ActiveStates = []int64{
		OrderStateCreated.Int(),
		OrderStateWaitingPayment.Int(),
		OrderStateBuyGoods.Int(),
	}

	// AllStates contains a list of order states
	AllStates = append([]int64{OrderStateDone.Int()}, ActiveStates...)
)

// OrderDoneState represents the different states of a completed order.
type OrderDoneState int64

// Constants defining various states of a completed order.
const (
	// OrderDoneStateSuccess represents the state when the order is completed successfully.
	OrderDoneStateSuccess OrderDoneState = iota
	// OrderDoneStateTimeout represents the state when the order has timed out.
	OrderDoneStateTimeout
	// OrderDoneStateCancel represents the state when the user has canceled the order.
	OrderDoneStateCancel
	// OrderDoneStatePurchaseFailed represents the state when the purchase has failed.
	OrderDoneStatePurchaseFailed
)

// String returns the string representation of the order done state.
func (s OrderDoneState) String() string {
	switch s {
	case 0:
		return "Success"
	case 1:
		return "Timeout"
	case 2:
		return "Cancel"
	case 3:
		return "PurchaseFailed"
	}

	return "Not found"
}

// Int returns the int representation of the order done state.
func (s OrderDoneState) Int() int64 {
	return int64(s)
}
