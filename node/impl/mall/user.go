package mall

import (
	"context"

	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"golang.org/x/xerrors"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
)

func (m *Mall) GetBalance(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	return m.LoadUserBalance(userID)
}

func (m *Mall) GetRechargeAddress(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)
	addr, err := m.GetRechargeAddressOfUser(userID)
	if err != nil {
		return "", err
	}

	if addr == "" {
		return m.TransactionMgr.AllocateTronAddress(userID)
	}

	return addr, nil
}

func (m *Mall) Withdraw(ctx context.Context, withdrawAddr, value string) error {
	userID := handler.GetID(ctx)

	return m.WithdrawManager.CreateWithdrawOrder(userID, withdrawAddr, value)
}

func (m *Mall) GetUserRechargeRecords(ctx context.Context, limit, offset int64) (*types.RechargeResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadRechargeRecordsByUser(userID, limit, offset)
}

func (m *Mall) GetUserWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) {
	userID := handler.GetID(ctx)

	return m.LoadWithdrawRecordsByUser(userID, limit, offset)
}

func (m *Mall) GetUserInstanceRecords(ctx context.Context, limit, offset int64) (*types.MyInstanceResponse, error) {
	userID := handler.GetID(ctx)
	k, s := m.getAccessKeys()
	instanceInfos, err := m.LoadMyInstancesInfo(userID, limit, offset)
	if err != nil {
		log.Errorf("LoadMyInstancesInfo err: %s", err.Error())
		return nil, xerrors.New(err.Error())
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

	return m.LoadInstanceDetailsInfo(userID, instanceID)
}
