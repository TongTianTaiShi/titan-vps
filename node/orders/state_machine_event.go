package orders

import (
	"golang.org/x/xerrors"
)

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

// PullAssetFatalError represents a fatal error in asset pulling
type PullAssetFatalError struct{ error }

// FormatError Format error
func (evt PullAssetFatalError) FormatError(xerrors.Printer) (next error) { return evt.error }

func (evt PullAssetFatalError) applyGlobal(state *OrderInfo) bool {
	log.Errorf("Fatal error on asset %s: %+v", state.Hash, evt.error)
	return true
}

// AssetForceState forces an asset state
type AssetForceState struct {
	State      OrderState
	Requester  string
	Details    string
	SeedNodeID string
}

func (evt AssetForceState) applyGlobal(state *OrderInfo) bool {
	state.State = evt.State
	return true
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

// PulledResult represents the result of node pulling
type PulledResult struct {
	BlocksCount int64
	Size        int64
}

func (evt PulledResult) apply(state *OrderInfo) {
}

func (evt PulledResult) Ignore() {
}

// PullRequestSent indicates that a pull request has been sent
type PullRequestSent struct{}

func (evt PullRequestSent) apply(state *OrderInfo) {
}

// AssetRePull re-pull the asset
type AssetRePull struct{}

func (evt AssetRePull) apply(state *OrderInfo) {
}

func (evt AssetRePull) Ignore() {
}

// PullSucceed indicates that a node has successfully pulled an asset
type PullSucceed struct{}

func (evt PullSucceed) apply(state *OrderInfo) {
}

func (evt PullSucceed) Ignore() {
}

// SkipStep skips the current step
type SkipStep struct{}

func (evt SkipStep) apply(state *OrderInfo) {}

// PullFailed indicates that a node has failed to pull an asset
type PullFailed struct{ error }

// FormatError Format error
func (evt PullFailed) FormatError(xerrors.Printer) (next error) { return evt.error }

func (evt PullFailed) apply(state *OrderInfo) {
}

func (evt PullFailed) Ignore() {
}

// SelectFailed  indicates that node selection has failed
type SelectFailed struct{ error }

// FormatError Format error
func (evt SelectFailed) FormatError(xerrors.Printer) (next error) { return evt.error }

func (evt SelectFailed) apply(state *OrderInfo) {
}
