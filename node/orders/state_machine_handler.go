package orders

import (
	"time"

	"github.com/LMF709268224/titan-vps/node/utils"
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

	height := m.getHeight()
	original, err := m.LoadUserBalance(info.User)
	if err != nil {
		return ctx.Send(OrderCancel{Height: height})
	}

	newValue, ok := utils.BigIntReduce(original, info.Value)
	if !ok {
		return ctx.Send(OrderCancel{Height: height})
	}

	err = m.UpdateUserBalance(info.User, newValue, original)
	if err != nil {
		return ctx.Send(OrderCancel{Height: height})
	}

	return ctx.Send(PaymentSucceed{PaymentInfo: info.PaymentInfo})

	// if info.PaymentInfo != nil {
	// 	if info.To == info.PaymentInfo.To {
	// 		orderAmount := new(big.Int)
	// 		orderAmount.SetString(info.Value, 10)

	// 		paymentAmount := new(big.Int)
	// 		paymentAmount.SetString(info.PaymentInfo.Value, 10)

	// 		if orderAmount.Cmp(paymentAmount) <= 0 {
	// 			return ctx.Send(PaymentSucceed{PaymentInfo: info.PaymentInfo})
	// 		}
	// 	}
	// }

	// return nil
}

// handleBuyGoods handles the order to buy goods
func (m *Manager) handleBuyGoods(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle buy goods: %s", info.OrderID)

	height := m.getHeight()

	// Buy Vps
	vInfo, err := m.LoadVpsInfo(info.VpsID)
	if err != nil {
		return ctx.Send(BuyFailed{Height: height, Msg: err.Error()})
	}

	rsp, err := m.createAliyunInstance(vInfo)
	if err != nil {
		return ctx.Send(BuyFailed{Height: height, Msg: err.Error()})
	}

	// Save To DB
	err = m.SaveVpsInstanceDevice(rsp)
	if err != nil {
		log.Errorf("SaveVpsInstanceDevice err:%s", err.Error())
	}

	return ctx.Send(BuySucceed{GoodsInfo: &GoodsInfo{ID: "vps_id", Password: "abc"}, Height: height})
}

// handleDone handles the order done
func (m *Manager) handleDone(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle done, %s, goods info:%v", info.OrderID, info.GoodsInfo)

	m.removeOrder(info.User)

	return nil
}
