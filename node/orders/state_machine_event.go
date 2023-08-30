package orders

type mutator interface {
	apply(state *OrderInfo)
}

// globalMutator is an event which can apply in every state
type globalMutator interface {
	// applyGlobal applies the event to the state. If if returns true,
	//  event processing should be interrupted
	applyGlobal(state *OrderInfo) bool
}

// Ignorable Ignorable
type Ignorable interface {
	Ignore()
}

// Global events

// OrderRestart restarts incomplete orders
type OrderRestart struct{}

func (evt OrderRestart) applyGlobal(state *OrderInfo) bool {
	return false
}

// PaymentResult represents the result of a user payment.
type PaymentResult struct{}

func (evt PaymentResult) apply(state *OrderInfo) {
}

// Ignore Ignorable
func (evt PaymentResult) Ignore() {
}

// CreateOrder represents the creation of a goods order.
type CreateOrder struct {
	*OrderInfo
}

func (evt CreateOrder) applyGlobal(state *OrderInfo) bool {
	state.OrderType = evt.OrderType
	state.State = evt.State
	state.OrderID = evt.OrderID
	state.User = evt.User
	state.Value = evt.Value
	state.DoneState = evt.DoneState
	state.VpsID = evt.VpsID
	state.EndTime = evt.EndTime

	return true
}

// WaitingPaymentSent indicates that the order is waiting for the user to make a payment.
type WaitingPaymentSent struct{}

func (evt WaitingPaymentSent) apply(state *OrderInfo) {}

// OrderTimeOut represents an order timeout event.
type OrderTimeOut struct{}

func (evt OrderTimeOut) apply(state *OrderInfo) {
	state.DoneState = OrderDoneStateTimeout
}

// OrderCancel represents an order cancellation event.
type OrderCancel struct{}

func (evt OrderCancel) apply(state *OrderInfo) {
	state.DoneState = OrderDoneStateCancel
}

// PaymentSucceed indicates a successful payment for the order.
type PaymentSucceed struct{}

func (evt PaymentSucceed) apply(state *OrderInfo) {
	// state.From = evt.From
	// state.TxHash = evt.TxHash
}

// BuySucceed represents a successful purchase event.
type BuySucceed struct {
	*GoodsInfo
}

func (evt BuySucceed) apply(state *OrderInfo) {
	state.GoodsInfo = evt.GoodsInfo
	state.DoneState = OrderDoneStateSuccess
}

// BuyFailed indicates a failed VPS purchase event.
type BuyFailed struct {
	Msg string
}

func (evt BuyFailed) apply(state *OrderInfo) {
	state.DoneState = OrderDoneStatePurchaseFailed
	state.Msg = evt.Msg
}
