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

// InfoUpdate update order info
type InfoUpdate struct {
	Size   int64
	Blocks int64
}

func (evt InfoUpdate) applyGlobal(state *OrderInfo) bool {
	return true
}

func (evt InfoUpdate) Ignore() {
}

// PaymentResult User payment result
type PaymentResult struct {
	*PaymentInfo
}

func (evt PaymentResult) apply(state *OrderInfo) {
	state.PaymentInfo = evt.PaymentInfo
}

func (evt PaymentResult) Ignore() {
}

// WaitingPaymentSent Waiting for user to pay
type WaitingPaymentSent struct{}

func (evt WaitingPaymentSent) apply(state *OrderInfo) {}

// OrderTimeOut order timeout
type OrderTimeOut struct{}

func (evt OrderTimeOut) apply(state *OrderInfo) {
	state.DoneState = Timeout
}

// OrderCancel cancel the order
type OrderCancel struct{}

func (evt OrderCancel) apply(state *OrderInfo) {
	state.DoneState = Cancel
}

// PaymentSucceed Order paid successfully
type PaymentSucceed struct{}

func (evt PaymentSucceed) apply(state *OrderInfo) {
}

func (evt PaymentSucceed) Ignore() {
}

// BuySucceed Successful purchase
type BuySucceed struct {
	*GoodsInfo
}

func (evt BuySucceed) apply(state *OrderInfo) {
	state.GoodsInfo = evt.GoodsInfo
	state.DoneState = Success
}
