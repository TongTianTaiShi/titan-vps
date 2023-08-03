package orders

// OrderState represents the state of an order in the process of being pulled.
type OrderState int64

// Constants defining various states of the order process.
const (
	// Created order
	Created OrderState = iota
	// WaitingPayment Waiting for user to payment order
	WaitingPayment
	// BuyGoods buy goods
	BuyGoods
	// Done the order done
	Done
)

// String returns the string representation of the order state.
func (s OrderState) String() string {
	switch s {
	case 0:
		return ""
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
	// PullingStates contains a list of order states that represent pulling.
	PullingStates = []int64{
		Created.Int(),
		WaitingPayment.Int(),
		BuyGoods.Int(),
	}

	// ActiveStates contains a list of order states that represent active.
	ActiveStates = append([]int64{Done.Int()}, PullingStates...)
)

// OrderDoneState represents the state of an order in the process of being pulled.
type OrderDoneState int64

// Constants defining various states of the done order process.
const (
	// Success order
	Success OrderDoneState = iota
	// Timeout timeout state
	Timeout
	// Cancel user cancel state
	Cancel
	// PurchaseFailed Purchase failed
	PurchaseFailed
)

// String returns the string representation of the order state.
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
