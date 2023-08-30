package mall

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/LMF709268224/titan-vps/node/utils"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
)

// GetAdminSignCode generates a sign code for an admin user.
func (m *Mall) GetAdminSignCode(ctx context.Context, userID string) (string, error) {
	exist, err := m.AdminExists(userID)
	if err != nil {
		return "", &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if !exist {
		return "", &api.ErrWeb{Code: terrors.NotAdministrator.Int(), Message: terrors.NotAdministrator.String()}
	}

	return m.UserMgr.GenerateSignCode(userID), nil
}

// LoginAdmin authenticates an admin user and generates a JWT token.
func (m *Mall) LoginAdmin(ctx context.Context, user *types.UserReq) (*types.LoginResponse, error) {
	userID := user.UserId

	code, err := m.UserMgr.GetSignCode(userID)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.NotFoundSignCode.Int(), Message: terrors.NotFoundSignCode.String()}
	}

	signature := user.Signature
	address, err := verifyEthMessage(code, signature)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	if strings.ToLower(userID) != strings.ToLower(address) {
		return nil, &api.ErrWeb{Code: terrors.UserMismatch.Int(), Message: fmt.Sprintf("%s,%s", userID, address)}
	}

	exist, err := m.AdminExists(address)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if !exist {
		return nil, &api.ErrWeb{Code: terrors.NotAdministrator.Int(), Message: terrors.NotAdministrator.String()}
	}

	p := types.JWTPayload{
		ID:        address,
		LoginType: int64(user.Type),
		Allow:     []auth.Permission{api.RoleAdmin},
	}

	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	rsp := &types.LoginResponse{}
	rsp.UserId = address
	rsp.Token = string(tk)

	return rsp, nil
}

// GetRechargeAddresses retrieves recharge addresses with pagination.
func (m *Mall) GetRechargeAddresses(ctx context.Context, limit, page int64) (*types.GetRechargeAddressResponse, error) {
	info, err := m.LoadRechargeAddresses(limit, page)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

// GetWithdrawalRecords retrieves withdrawal records with optional filtering.
func (m *Mall) GetWithdrawalRecords(ctx context.Context, req *types.GetWithdrawRequest) (*types.GetWithdrawResponse, error) {
	statuses := make([]types.WithdrawState, 0)
	if req.State == "" {
		statuses = []types.WithdrawState{types.WithdrawCreate, types.WithdrawDone, types.WithdrawRefund}
	} else {
		s2, err := strconv.Atoi(req.State)
		if err != nil {
			return nil, &api.ErrWeb{Code: terrors.ParametersWrong.Int(), Message: fmt.Sprintf("state is %s , err:%s", req.State, err.Error())}
		}

		statuses = []types.WithdrawState{types.WithdrawState(s2)}
	}

	info, err := m.LoadWithdrawRecords(req.Limit, req.Offset, statuses, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

// ApproveUserWithdrawal approves a user withdrawal request.
func (m *Mall) ApproveUserWithdrawal(ctx context.Context, orderID, withdrawHash string) error {
	userID := handler.GetID(ctx)

	info, err := m.LoadWithdrawRecord(orderID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if info.State != types.WithdrawCreate {
		return &api.ErrWeb{Code: terrors.StatusNotEditable.Int(), Message: string(info.State)}
	}

	info.WithdrawHash = withdrawHash
	info.Executor = userID

	err = m.UpdateWithdrawRecord(info, types.WithdrawDone)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return nil
}

// RejectUserWithdrawal rejects a user withdrawal request.
func (m *Mall) RejectUserWithdrawal(ctx context.Context, orderID string) error {
	userID := handler.GetID(ctx)

	info, err := m.LoadWithdrawRecord(orderID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if info.State != types.WithdrawCreate {
		return &api.ErrWeb{Code: terrors.StatusNotEditable.Int(), Message: string(info.State)}
	}

	info.Executor = userID

	err = m.UpdateWithdrawRecord(info, types.WithdrawRefund)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	original, err := m.LoadUserBalance(info.UserID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	newValue, err := utils.AddBigInt(original, info.Value)
	if err != nil {
		return err
	}

	err = m.UpdateUserBalance(info.UserID, newValue, original)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return nil
}

// AddAdminUser adds an admin user with a userID and nickname.
func (m *Mall) AddAdminUser(ctx context.Context, userID, nickName string) error {
	err := m.SaveAdminInfo(userID, nickName)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return nil
}

// SupplementRechargeOrder supplements a recharge order.
func (m *Mall) SupplementRechargeOrder(ctx context.Context, hash string) error {
	return m.TransactionMgr.SupplementOrder(hash)
}

// RefundInstance is a method that handles refund request for a specific instance
func (m *Mall) RefundInstance(ctx context.Context, instanceID string) (int64, error) {
	accessKeyID, accessKeySecret := m.getAliAccessKeys()

	return aliyun.RefundInstance(accessKeyID, accessKeySecret, instanceID)
}

// InquiryPriceRefundInstance is a method that inquires the price of refunding a specific instance
func (m *Mall) InquiryPriceRefundInstance(ctx context.Context, instanceID string) (float64, error) {
	accessKeyID, accessKeySecret := m.getAliAccessKeys()

	return aliyun.InquiryPriceRefundInstance(accessKeyID, accessKeySecret, instanceID)
}

// GetInstanceRecords is a method that retrieves the records of instances
func (m *Mall) GetInstanceRecords(ctx context.Context, limit, page int64) (*types.GetInstanceResponse, error) {
	accessKeyID, accessKeySecret := m.getAliAccessKeys()
	instanceInfos, err := m.LoadInstancesInfo(limit, page)
	if err != nil {
		log.Errorf("LoadMyInstancesInfo err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	for _, instanceInfo := range instanceInfos.List {
		var instanceIds []string
		instanceIds = append(instanceIds, instanceInfo.InstanceId)

		rsp, err := aliyun.DescribeInstanceStatus(instanceInfo.RegionId, accessKeyID, accessKeySecret, instanceIds)
		if err != nil {
			log.Errorf("DescribeInstanceStatus err: %s", *err.Message)
			continue
		}

		if len(rsp.Body.InstanceStatuses.InstanceStatus) == 0 {
			continue
		}

		instanceInfo.State = *rsp.Body.InstanceStatuses.InstanceStatus[0].Status
		instanceInfo.Renew = ""
		if instanceInfo.State == "Stopped" {
			continue
		}

		renewInfo := types.SetRenewOrderReq{
			RegionID:   instanceInfo.RegionId,
			InstanceId: instanceInfo.InstanceId,
		}

		status, errEk := m.GetRenewInstance(ctx, renewInfo)
		if errEk != nil {
			log.Errorf("GetRenewInstance err: %s", errEk.Error())
			continue
		}
		instanceInfo.Renew = status
		instanceExpiredTime, err := aliyun.DescribeInstances(instanceInfo.RegionId, accessKeyID, accessKeySecret, instanceIds)
		if err != nil {
			log.Errorf("DescribeInstances err: %s", *err.Message)
			continue
		}
		if len(instanceExpiredTime.Body.Instances.Instance) > 0 {
			instanceInfo.ExpiredTime = *instanceExpiredTime.Body.Instances.Instance[0].ExpiredTime
		}
	}

	return instanceInfos, nil
}
