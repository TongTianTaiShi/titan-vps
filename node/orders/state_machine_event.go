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

// PaymentResult User payment result
type PaymentResult struct {
	*PaymentInfo
}

func (evt PaymentResult) apply(state *OrderInfo) {
	state.PaymentInfo = evt.PaymentInfo
}

// Ignore Ignorable
func (evt PaymentResult) Ignore() {
}

// CreateOrder create a goods order
type CreateOrder struct {
	*OrderInfo
}

func (evt CreateOrder) applyGlobal(state *OrderInfo) bool {
	state.CreatedHeight = evt.CreatedHeight
	state.To = evt.To
	state.State = evt.State
	state.OrderID = evt.OrderID
	state.From = evt.From
	state.User = evt.User
	state.Value = evt.Value
	state.DoneState = evt.DoneState
	state.DoneHeight = evt.DoneHeight
	state.VpsID = evt.VpsID

	return true
}

// WaitingPaymentSent Waiting for user to pay
type WaitingPaymentSent struct{}

func (evt WaitingPaymentSent) apply(state *OrderInfo) {}

// OrderTimeOut order timeout
type OrderTimeOut struct {
	Height int64
}

func (evt OrderTimeOut) apply(state *OrderInfo) {
	state.DoneState = Timeout
	state.DoneHeight = evt.Height
}

// OrderCancel cancel the order
type OrderCancel struct {
	Height int64
}

func (evt OrderCancel) apply(state *OrderInfo) {
	state.DoneState = Cancel
	state.DoneHeight = evt.Height
}

// PaymentSucceed Order paid successfully
type PaymentSucceed struct {
	*PaymentInfo
}

func (evt PaymentSucceed) apply(state *OrderInfo) {
	state.From = evt.From
	state.TxHash = evt.TxHash
}

// BuySucceed Successful purchase
type BuySucceed struct {
	*GoodsInfo
	Height int64
}

func (evt BuySucceed) apply(state *OrderInfo) {
	state.GoodsInfo = evt.GoodsInfo
	state.DoneState = Success
	state.DoneHeight = evt.Height
}

// BuyFailed buy vps failed
type BuyFailed struct {
	Height int64
	Msg    string
}

func (evt BuyFailed) apply(state *OrderInfo) {
	state.DoneState = PurchaseFailed
	state.DoneHeight = evt.Height
	state.Msg = evt.Msg
}
