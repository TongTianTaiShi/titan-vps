package orders

// OrderState represents the state of an asset in the process of being pulled.
type OrderState int64

// Constants defining various states of the asset pulling process.
const (
	// Created select first candidate to pull seed asset
	Created OrderState = iota
	// WaitingPayment Waiting for candidate nodes to pull seed asset
	WaitingPayment
	// BuyGoods Initialize user upload preparation
	BuyGoods
	// Done Waiting for user to upload asset to candidate node
	Done
	// BuyGoodsFailed Unable to select candidate nodes or failed to pull asset
	BuyGoodsFailed
)

// String returns the string representation of the AssetState.
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
	case 4:
		return "BuyGoodsFailed"
	}

	return "Not found"
}

// Int returns the int representation of the AssetState.
func (s OrderState) Int() int64 {
	return int64(s)
}

var (
	// FailedStates contains a list of asset pull states that represent failures.
	FailedStates = []int64{
		BuyGoodsFailed.Int(),
	}

	// PullingStates contains a list of asset pull states that represent pulling.
	PullingStates = []int64{
		Created.Int(),
		WaitingPayment.Int(),
		BuyGoods.Int(),
	}

	// ActiveStates contains a list of asset pull states that represent active.
	ActiveStates = append(append([]int64{Done.Int()}, FailedStates...), PullingStates...)
)
