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

// handleCreated handles the order create
func (m *Manager) handleCreated(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle order created , %s", info.OrderID)

	return ctx.Send(WaitingPaymentSent{})
}

// handleWaitingPayment handles the order wait for user payment
func (m *Manager) handleWaitingPayment(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle wait payment, %s , info : %v", info.OrderID, info.PaymentInfo)

	if info.PaymentInfo != nil {
		if info.To == info.PaymentInfo.To && info.Value <= info.PaymentInfo.Value {
			return ctx.Send(PaymentSucceed{})
		}
	}

	return nil
}

// handleBuyGoods handles the order to buy goods
func (m *Manager) handleBuyGoods(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle buy goods: %s", info.OrderID)

	// Buy Vps
	vInfo, err := m.LoadVpsInfo(info.VpsID)
	if err != nil {
		return err
	}

	_, err = m.createAliyunInstance(vInfo)
	// Save To DB

	height := m.filecoinMgr.GetHeight()

	return ctx.Send(BuySucceed{GoodsInfo: &GoodsInfo{ID: "vps_id", Password: "abc"}, Height: height})
}

// handleDone handles the order done
func (m *Manager) handleDone(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle done, %s, goods info:%v", info.OrderID, info.GoodsInfo)

	m.revertPayeeAddress(info.To)
	m.removeOrder(info.From)

	return nil
}
