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

// OrderRestart restarts asset pulling
type OrderRestart struct{}

func (evt OrderRestart) applyGlobal(state *OrderInfo) bool {
	return false
}

// InfoUpdate update asset info
type InfoUpdate struct {
	Size   int64
	Blocks int64
}

func (evt InfoUpdate) applyGlobal(state *OrderInfo) bool {
	return true
}

func (evt InfoUpdate) Ignore() {
}

// PaymentResult represents the result of node pulling
type PaymentResult struct {
	Info *PaymentInfo
}

func (evt PaymentResult) apply(state *OrderInfo) {
	state.PaymentInfo = evt.Info
}

func (evt PaymentResult) Ignore() {
}

// WaitingPaymentSent indicates that a pull request has been sent
type WaitingPaymentSent struct{}

func (evt WaitingPaymentSent) apply(state *OrderInfo) {}

// OrderTimeOut indicates that a pull request has been sent
type OrderTimeOut struct{}

func (evt OrderTimeOut) apply(state *OrderInfo) {}

// OrderCancel indicates that a pull request has been sent
type OrderCancel struct{}

func (evt OrderCancel) apply(state *OrderInfo) {}

// PaymentSucceed indicates that a node has successfully pulled an asset
type PaymentSucceed struct{}

func (evt PaymentSucceed) apply(state *OrderInfo) {
}

func (evt PaymentSucceed) Ignore() {
}

// BuySucceed skips the current step
type BuySucceed struct{}

func (evt BuySucceed) apply(state *OrderInfo) {}
