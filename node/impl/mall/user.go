package mall

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/lib/aliyun"

	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
)

// GetBalance Get user balance
func (m *Mall) GetBalance(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	balance, err := m.LoadUserBalance(userID)
	if err != nil {
		return balance, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return balance, nil
}

// GetRechargeAddress  Get user recharge address
func (m *Mall) GetRechargeAddress(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	address, err := m.LoadRechargeAddressOfUser(userID)
	if err != nil {
		return address, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return address, nil
}

// Withdraw user Withdraw
func (m *Mall) Withdraw(ctx context.Context, withdrawAddr, value string) error {
	userID := handler.GetID(ctx)

	return m.WithdrawManager.CreateWithdrawOrder(userID, withdrawAddr, value)
}

func (m *Mall) GetUserRechargeRecords(ctx context.Context, limit, offset int64) (*types.RechargeResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadRechargeRecordsByUser(userID, limit, offset)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

func (m *Mall) GetUserWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadWithdrawRecordsByUser(userID, limit, offset)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

func (m *Mall) GetUserInstanceRecords(ctx context.Context, limit, offset int64) (*types.MyInstanceResponse, error) {
	userID := handler.GetID(ctx)
	k, s := m.getAccessKeys()
	instanceInfos, err := m.LoadMyInstancesInfo(userID, limit, offset)
	if err != nil {
		log.Errorf("LoadMyInstancesInfo err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}
	for _, instanceInfo := range instanceInfos.List {
		var instanceIds []string
		instanceIds = append(instanceIds, instanceInfo.InstanceId)
		rsp, err := aliyun.DescribeInstanceStatus(instanceInfo.Location, k, s, instanceIds)
		if err != nil {
			log.Errorf("DescribeInstanceStatus err: %s", err.Error())
		}
		instanceInfo.State = *rsp.Body.InstanceStatuses.InstanceStatus[0].Status
	}

	return instanceInfos, nil
}

func (m *Mall) GetInstanceDetailsInfo(ctx context.Context, instanceID string) (*types.InstanceDetails, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadInstanceDetailsInfo(userID, instanceID)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}
