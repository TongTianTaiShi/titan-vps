package orders

// OrderState represents the state of an asset in the process of being pulled.
type OrderState string

// Constants defining various states of the asset pulling process.
const (
	// Created select first candidate to pull seed asset
	Created OrderState = ""
	// WaitingPayment Waiting for candidate nodes to pull seed asset
	WaitingPayment OrderState = "WaitingPayment"
	// BuyGoods Initialize user upload preparation
	BuyGoods OrderState = "BuyGoods"
	// Done Waiting for user to upload asset to candidate node
	Done OrderState = "Done"
	// PaymentFailed Unable to select candidate nodes or failed to pull seed asset
	PaymentFailed OrderState = "PaymentFailed"
	// BuyGoodsFailed Unable to select candidate nodes or failed to pull asset
	BuyGoodsFailed OrderState = "BuyGoodsFailed"
	// Remove remove
	Remove OrderState = "Remove"
)

// String returns the string representation of the AssetState.
func (s OrderState) String() string {
	return string(s)
}

var (
	// FailedStates contains a list of asset pull states that represent failures.
	FailedStates = []string{
		PaymentFailed.String(),
		BuyGoodsFailed.String(),
	}

	// PullingStates contains a list of asset pull states that represent pulling.
	PullingStates = []string{
		Created.String(),
		WaitingPayment.String(),
		BuyGoods.String(),
	}

	// ActiveStates contains a list of asset pull states that represent active.
	ActiveStates = append(append([]string{Done.String()}, FailedStates...), PullingStates...)
)
