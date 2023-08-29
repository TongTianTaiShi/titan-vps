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

// handleOrderCreated handles the order creati
func (m *Manager) handleOrderCreated(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle order created , %s", info.OrderID)

	return ctx.Send(WaitingPaymentSent{})
}

// handleWaitingForPayment handles the order waiting for user payment
func (m *Manager) handleWaitingForPayment(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle wait payment, %s ", info.OrderID)

	original, err := m.LoadUserBalance(info.User)
	if err != nil {
		log.Errorf("handleWaitingForPayment LoadUserBalance err:%s", err.Error())
		return nil
	}

	newValue, err := utils.ReduceBigInt(original, info.Value)
	if err != nil {
		log.Errorf("handleWaitingForPayment BigIntReduce err:%s", err.Error())
		return nil
	}

	err = m.UpdateUserBalance(info.User, newValue, original)
	if err != nil {
		log.Errorf("handleWaitingForPayment UpdateUserBalance err:%s", err.Error())
		return nil
	}

	return ctx.Send(PaymentSucceed{})
}

// handleBuyGoods handles the order to buy goods
func (m *Manager) handleBuyGoods(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle buy goods: %s", info.OrderID)

	// Buy Vps
	vInfo, err := m.LoadInstanceInfoByID(info.VpsID)
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
		createInfo := &types.CreateInstanceReq{
			RegionID:           vInfo.RegionId,
			InstanceType:       vInfo.InstanceType,
			ImageID:            vInfo.ImageID,
			SecurityGroupID:    vInfo.SecurityGroupId,
			PeriodUnit:         vInfo.PeriodUnit,
			Period:             vInfo.Period,
			DryRun:             vInfo.DryRun,
			InternetChargeType: vInfo.InternetChargeType,
			SystemDiskSize:     vInfo.SystemDiskSize,
			SystemDiskCategory: vInfo.SystemDiskCategory,
			BandwidthOut:       vInfo.BandwidthOut,
			DataDisk:           vInfo.DataDisk,
		}

		result, err := m.vpsMgr.CreateAliYunInstance(vInfo.OrderID, createInfo)
		if err != nil {
			return ctx.Send(BuyFailed{Msg: err.Error()})
		}
		vInfo.InstanceId = result.InstanceID
	} else if info.OrderType == int64(types.RenewVPS) {
		err = m.vpsMgr.RenewInstance(&types.RenewInstanceRequest{
			RegionId:   vInfo.RegionId,
			InstanceId: vInfo.InstanceId,
			PeriodUnit: vInfo.PeriodUnit,
			Period:     vInfo.Period,
		})
		if err != nil {
			return ctx.Send(BuyFailed{Msg: err.Error()})
		}
	}

	// if auto renew
	if vInfo.AutoRenew == 1 {
		renewReq := types.SetRenewOrderReq{
			RegionID:   vInfo.RegionId,
			InstanceId: vInfo.InstanceId,
			PeriodUnit: vInfo.PeriodUnit,
			Period:     vInfo.Period,
			Renew:      1,
		}
		err = m.vpsMgr.ModifyInstanceRenew(&renewReq)
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

// handleOrderDone handles the order completion
func (m *Manager) handleOrderDone(ctx statemachine.Context, info OrderInfo) error {
	log.Debugf("handle done, %s, goods info:%v", info.OrderID, info.GoodsInfo)

	if info.DoneState == OrderDoneStatePurchaseFailed {
		original, err := m.LoadUserBalance(info.User)
		if err != nil {
			log.Errorf("handleOrderDone LoadUserBalance err:%s", err.Error())
			return nil
		}

		newValue, err := utils.AddBigInt(original, info.Value)
		if err != nil {
			log.Errorf("handleOrderDone BigIntAdd err:%s", err.Error())
			return nil
		}

		err = m.UpdateUserBalance(info.User, newValue, original)
		if err != nil {
			log.Errorf("handleOrderDone UpdateUserBalance err:%s", err.Error())
			return nil
		}
	}

	m.removeOrder(info.OrderID.String())

	return nil
}
