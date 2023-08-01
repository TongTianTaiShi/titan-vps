package orders

import (
	"time"

	"github.com/filecoin-project/go-statemachine"
)

var (
	// MinRetryTime defines the minimum time duration between retries
	MinRetryTime = 1 * time.Minute

	// MaxRetryCount defines the maximum number of retries allowed
	MaxRetryCount = 3
)

// failedCoolDown is called when a retry needs to be attempted and waits for the specified time duration
func failedCoolDown(ctx statemachine.Context, info OrderInfo) error {
	retryStart := time.Now().Add(MinRetryTime)
	if time.Now().Before(retryStart) {
		log.Debugf("%s(%s), waiting %s before retrying", info.State, info.Hash, time.Until(retryStart))
		select {
		case <-time.After(time.Until(retryStart)):
		case <-ctx.Context().Done():
			return ctx.Context().Err()
		}
	}

	return nil
}

// handleCreated handles the selection of seed nodes for asset pull
func (m *Manager) handleCreated(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle select seed: %s", info.Hash)

	return ctx.Send(PullRequestSent{})
}

// handleWaitingPayment handles the asset pulling process of seed nodes
func (m *Manager) handleWaitingPayment(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle seed pulling, %s", info.Hash)

	return nil
}

// handleBuyGoods handles the upload init
func (m *Manager) handleBuyGoods(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle upload init: %s", info.Hash)

	return ctx.Send(PullRequestSent{})
}

// handleDone handles the asset upload process of seed nodes
func (m *Manager) handleDone(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle seed upload, %s", info.Hash)

	return nil
}

// handleCandidatesSelect handles the selection of candidate nodes for asset pull
func (m *Manager) handleCandidatesSelect(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle candidates select, %s", info.Hash)

	return ctx.Send(PullRequestSent{})
}

// handleCandidatesPulling handles the asset pulling process of candidate nodes
func (m *Manager) handleCandidatesPulling(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle candidates pulling, %s", info.Hash)

	return nil
}

// handleEdgesSelect handles the selection of edge nodes for asset pull
func (m *Manager) handleEdgesSelect(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle edges select , %s", info.Hash)

	return ctx.Send(PullRequestSent{})
}

// handleEdgesPulling handles the asset pulling process of edge nodes
func (m *Manager) handleEdgesPulling(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle edges pulling, %s", info.Hash)

	return nil
}

// handleServicing asset pull completed and in service status
func (m *Manager) handleServicing(ctx statemachine.Context, info OrderInfo) error {
	log.Infof("handle servicing: %s", info.Hash)
	// remove fail replicas
	return nil
}

// handlePullsFailed handles the failed state of asset pulling and retries if necessary
func (m *Manager) handlePullsFailed(ctx statemachine.Context, info OrderInfo) error {
	return ctx.Send(AssetRePull{})
}

func (m *Manager) handleUploadFailed(ctx statemachine.Context, info OrderInfo) error {
	log.Infof("handle upload fail: %s", info.Hash)

	return nil
}

func (m *Manager) handleRemove(ctx statemachine.Context, info OrderInfo) error {
	log.Infof("handle remove: %s", info.Hash)

	return nil
}
