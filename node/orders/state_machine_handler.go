package orders

import (
	"encoding/json"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
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
	log.Debugf("handle wait payment, %s ", info.OrderID)

	original, err := m.LoadUserBalance(info.User)
	if err != nil {
		log.Errorf("WaitingPayment LoadUserBalance err:%s", err.Error())
		return nil
	}

	newValue, err := utils.BigIntReduce(original, info.Value)
	if err != nil {
		log.Errorf("WaitingPayment BigIntReduce err:%s", err.Error())
		return nil
	}

	err = m.UpdateUserBalance(info.User, newValue, original)
	if err != nil {
		log.Errorf("WaitingPayment UpdateUserBalance err:%s", err.Error())
		return nil
	}

	return ctx.Send(PaymentSucceed{})
}

// handleBuyGoods handles the order to buy goods
func (m *Manager) handleBuyGoods(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle buy goods: %s", info.OrderID)

	// Buy Vps
	vInfo, err := m.LoadVpsInfo(info.VpsID)
	if err != nil {
		return ctx.Send(BuyFailed{Msg: err.Error()})
	}
	vInfo.UserID = info.User
	vInfo.OrderID = info.OrderID.String()
	if vInfo.DataDiskString != "" {
		if err := json.Unmarshal([]byte(vInfo.DataDiskString), &vInfo.DataDisk); err != nil {
			return ctx.Send(BuyFailed{Msg: err.Error()})
		}
	}

	if info.OrderType == int64(types.BuyVPS) {
		result, err := m.vMgr.CreateAliYunInstance(vInfo)
		if err != nil {
			return ctx.Send(BuyFailed{Msg: err.Error()})
		}
		vInfo.InstanceId = result.InstanceID
	} else if info.OrderType == int64(types.RenewVPS) {
		err = m.vMgr.RenewInstance(&types.RenewInstanceRequest{
			RegionId:   vInfo.RegionId,
			InstanceId: vInfo.InstanceId,
			PeriodUnit: vInfo.PeriodUnit,
			Period:     vInfo.Period,
		})
		if err != nil {
			return ctx.Send(BuyFailed{Msg: err.Error()})
		}
	}
	if vInfo.Renew == 1 {
		renewReq := types.SetRenewOrderReq{
			RegionID:   vInfo.RegionId,
			InstanceId: vInfo.InstanceId,
			PeriodUnit: vInfo.PeriodUnit,
			Period:     vInfo.Period,
			Renew:      1,
		}
		err = m.vMgr.ModifyInstanceRenew(&renewReq)
		if err != nil {
			log.Errorf("ModifyInstanceRenew err: %v", err)
		}
	}
	//// Save To DB
	//err = m.SaveVpsInstanceDevice(rsp)
	//if err != nil {
	//	log.Errorf("SaveVpsInstanceDevice err:%s", err.Error())
	//}

	return ctx.Send(BuySucceed{GoodsInfo: &GoodsInfo{ID: "vps_id", Password: "abc"}})
}

// handleDone handles the order done
func (m *Manager) handleDone(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle done, %s, goods info:%v", info.OrderID, info.GoodsInfo)

	if info.DoneState == PurchaseFailed {
		original, err := m.LoadUserBalance(info.User)
		if err != nil {
			log.Errorf("handleDone LoadUserBalance err:%s", err.Error())
			return nil
		}

		newValue, err := utils.BigIntAdd(original, info.Value)
		if err != nil {
			log.Errorf("handleDone BigIntAdd err:%s", err.Error())
			return nil
		}

		err = m.UpdateUserBalance(info.User, newValue, original)
		if err != nil {
			log.Errorf("handleDone UpdateUserBalance err:%s", err.Error())
			return nil
		}
	}

	m.removeOrder(info.OrderID.String())

	return nil
}
