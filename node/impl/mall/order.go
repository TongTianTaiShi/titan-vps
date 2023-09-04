package mall

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
)

const decimal = 1000000

func countEndDate(unit string, period int) time.Time {
	tt := time.Now()

	switch unit {
	case "Week":
		tt = tt.AddDate(0, 0, 7*period)
	case "Month":
		tt = tt.AddDate(0, period, 0)
	case "Year":
		tt = tt.AddDate(period, 0, 0)
	}

	return tt
}

// CreateOrder creates a new order.
func (m *Mall) CreateOrder(ctx context.Context, req types.CreateOrderReq) (string, error) {
	userID := handler.GetID(ctx)

	instanceDetails := &types.InstanceDetails{
		RegionId:           req.RegionId,
		InstanceType:       req.InstanceType,
		ImageID:            req.ImageID,
		SecurityGroupId:    req.SecurityGroupID,
		PeriodUnit:         req.PeriodUnit,
		Period:             req.Period,
		InternetChargeType: req.InternetChargeType,
		SystemDiskSize:     req.SystemDiskSize,
		SystemDiskCategory: req.SystemDiskCategory,
		BandwidthOut:       req.InternetMaxBandwidthOut,
		UserID:             userID,
		InstanceChargeType: req.InstanceChargeType,
		BandwidthIn:        req.InternetMaxBandwidthIn,
		AutoRenew:          req.Renew,
		State:              "Pending",
	}

	// Marshal DataDisk if it's not empty
	if len(req.DataDisk) > 0 {
		dataDisk, err := json.Marshal(req.DataDisk)
		if err != nil {
			log.Errorf("Marshal DataDisk:%v", err)
			return "", &api.ErrWeb{Code: terrors.ParametersWrong.Int(), Message: err.Error()}
		}
		instanceDetails.DataDiskString = string(dataDisk)
	}

	priceReq := &types.DescribePriceReq{
		RegionId:                     req.RegionId,
		InstanceType:                 req.InstanceType,
		PriceUnit:                    req.PeriodUnit,
		Period:                       req.Period,
		Amount:                       req.Amount,
		InternetChargeType:           req.InternetChargeType,
		ImageID:                      req.ImageID,
		InternetMaxBandwidthOut:      req.InternetMaxBandwidthOut,
		SystemDiskCategory:           req.SystemDiskCategory,
		SystemDiskSize:               req.SystemDiskSize,
		DescribePriceRequestDataDisk: req.DataDisk,
	}

	priceInfo, err := m.DescribePrice(ctx, priceReq)
	if err != nil {
		log.Errorf("DescribePrice:%v", err)
		return "", &api.ErrWeb{Code: terrors.DescribePriceError.Int(), Message: err.Error()}
	}
	// Calculate new balance
	newBalanceString := strconv.FormatFloat(math.Ceil(float64(priceInfo.USDPrice)*decimal), 'f', 0, 64)

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)
	instanceDetails.OrderID = orderID
	instanceDetails.Value = newBalanceString

	id, err := m.SaveInstanceInfoOfUser(instanceDetails)
	if err != nil {
		log.Errorf("SaveVpsInstance:%v", err)
		return "", err
	}

	endDate := countEndDate(req.PeriodUnit, int(req.Period))

	// Create an order record
	info := &types.OrderRecord{
		VpsID:     id,
		OrderID:   orderID,
		UserID:    userID,
		Value:     newBalanceString,
		OrderType: types.BuyVPS,
		CycleTime: fmt.Sprintf("%s - %s", time.Now().Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")),
	}

	err = m.OrderMgr.CreatedOrder(info)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

// RenewOrder renews an existing order.
func (m *Mall) RenewOrder(ctx context.Context, renewReq types.RenewOrderReq) (string, error) {
	userID := handler.GetID(ctx)

	req, err := m.LoadUserInstanceInfoByInstanceID(renewReq.InstanceId)
	if err != nil {
		return "", &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if len(req.DataDisk) > 0 {
		dataDisk, err := json.Marshal(req.DataDisk)
		if err != nil {
			log.Errorf("Marshal DataDisk:%v", err)
			return "", &api.ErrWeb{Code: terrors.ParametersWrong.Int(), Message: err.Error()}
		}
		req.DataDiskString = string(dataDisk)
	}

	priceReq := &types.DescribePriceReq{
		RegionId:                     req.RegionId,
		InstanceType:                 req.InstanceType,
		PriceUnit:                    renewReq.PeriodUnit,
		Period:                       renewReq.Period,
		Amount:                       1,
		InternetChargeType:           req.InternetChargeType,
		ImageID:                      req.ImageID,
		InternetMaxBandwidthOut:      req.BandwidthOut,
		SystemDiskCategory:           req.SystemDiskCategory,
		SystemDiskSize:               req.SystemDiskSize,
		DescribePriceRequestDataDisk: req.DataDisk,
	}

	priceInfo, err := m.DescribePrice(ctx, priceReq)
	if err != nil {
		log.Errorf("DescribePrice:%v", err)
		return "", &api.ErrWeb{Code: terrors.DescribePriceError.Int(), Message: err.Error()}
	}

	newBalanceString := strconv.FormatFloat(math.Ceil(float64(priceInfo.USDPrice)*decimal), 'f', 0, 64)

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)
	// req.OrderID = orderID
	req.Value = newBalanceString
	req.PeriodUnit = renewReq.PeriodUnit
	req.Period = renewReq.Period
	req.AutoRenew = renewReq.Renew

	err = m.RenewVpsInstance(req)
	if err != nil {
		log.Errorf("SaveVpsInstance:%v", err)
		return "", err
	}

	endDate := countEndDate(req.PeriodUnit, int(req.Period))

	eTime, err := time.Parse("2006-01-02T15:04Z", req.ExpiredTime)

	info := &types.OrderRecord{
		VpsID:     req.ID,
		OrderID:   orderID,
		UserID:    userID,
		Value:     newBalanceString,
		OrderType: types.RenewVPS,
		CycleTime: fmt.Sprintf("%s - %s", eTime.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")),
	}

	err = m.OrderMgr.CreatedOrder(info)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

// GetUseWaitingPaymentOrders retrieves user's unpaid orders with pagination.
func (m *Mall) GetUseWaitingPaymentOrders(ctx context.Context, limit, page int64) (*types.OrderRecordResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadOrderRecordByUserUndone(userID, limit, page, m.OrderMgr.GetOrderTimeoutDurationMinutes())
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

// GetUserOrderRecords retrieves user's order records with pagination.
func (m *Mall) GetUserOrderRecords(ctx context.Context, limit, page int64) (*types.OrderRecordResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadOrderRecordsByUser(userID, limit, page, m.OrderMgr.GetOrderTimeoutDurationMinutes())
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

// CancelUserOrder cancels a user's order.
func (m *Mall) CancelUserOrder(ctx context.Context, orderID string) error {
	userID := handler.GetID(ctx)
	return m.OrderMgr.CancelOrder(orderID, userID)
}

// PaymentUserOrder marks a user's order as paid.
func (m *Mall) PaymentUserOrder(ctx context.Context, orderID string) error {
	userID := handler.GetID(ctx)
	return m.OrderMgr.PaymentCompleted(orderID, userID)
}
