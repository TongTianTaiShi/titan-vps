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
		log.Debugf("%s(%s), waiting %s before retrying", info.State, info.OrderID, time.Until(retryStart))
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
	log.Debugf("handle select seed: %s", info.OrderID)

	return ctx.Send(WaitingPaymentSent{})
}

// handleWaitingPayment handles the asset pulling process of seed nodes
func (m *Manager) handleWaitingPayment(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle seed pulling, %s", info.OrderID)

	return nil
}

// handleBuyGoods handles the upload init
func (m *Manager) handleBuyGoods(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle upload init: %s", info.OrderID)

	return ctx.Send(WaitingPaymentSent{})
}

// handleDone handles the asset upload process of seed nodes
func (m *Manager) handleDone(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle seed upload, %s", info.OrderID)

	m.revertPayeeAddress(info.To)
	m.removeOrder(info.From)

	return nil
}
