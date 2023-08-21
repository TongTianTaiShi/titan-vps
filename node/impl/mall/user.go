package mall

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"

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
func (m *Mall) GetInstanceDefaultInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) (*types.InstanceTypeResponse, error) {
	req.Offset = req.Limit * (req.Page - 1)
	return m.LoadInstanceDefaultInfo(req)
}

func (m *Mall) GetSignCode(ctx context.Context, userID string) (string, error) {
	return m.UserMgr.GenerateSignCode(userID), nil
}

func (m *Mall) Login(ctx context.Context, user *types.UserReq) (*types.LoginResponse, error) {
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

	p := types.JWTPayload{
		ID:        address,
		LoginType: int64(user.Type),
		Allow:     []auth.Permission{api.RoleUser},
	}
	rsp := &types.LoginResponse{}
	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return rsp, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}
	rsp.UserId = address
	rsp.Token = string(tk)

	err = m.initUser(address)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (m *Mall) initUser(userID string) error {
	exist, err := m.UserExists(userID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if exist {
		return nil
	}

	// init recharge address
	addr, err := m.LoadRechargeAddressOfUser(userID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if addr == "" {
		_, err = m.TransactionMgr.AllocateTronAddress(userID)
		if err != nil {
			return err
		}
	}

	err = m.SaveUserInfo(&types.UserInfo{UserID: userID, Balance: "0"})
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return nil
}

func (m *Mall) Logout(ctx context.Context, user *types.UserReq) error {
	userID := handler.GetID(ctx)
	log.Warnf("user id : %s", userID)
	// delete(m.UserMgr.User, user.UserId)
	return nil
}