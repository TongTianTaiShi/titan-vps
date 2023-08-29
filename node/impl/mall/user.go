package mall

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/lib/trxbridge"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"

	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/LMF709268224/titan-vps/node/utils"
)

// GetBalance retrieves user balance.
func (m *Mall) GetBalance(ctx context.Context) (*types.UserInfo, error) {
	userID := handler.GetID(ctx)

	uInfo := &types.UserInfo{UserID: userID}

	balance, err := m.LoadUserBalance(userID)
	if err != nil {
		return uInfo, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	uInfo.Balance = balance

	list, err := m.LoadWithdrawRecordsByUserAndState(userID, types.WithdrawCreate)
	if err != nil {
		return uInfo, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	lockBalance := "0"
	for _, info := range list {
		b, err := utils.AddBigInt(info.Value, lockBalance)
		if err == nil {
			lockBalance = b
		}
	}

	uInfo.LockedBalance = lockBalance

	return uInfo, nil
}

// GetRechargeAddress retrieves user recharge address.
func (m *Mall) GetRechargeAddress(ctx context.Context) (string, error) {
	userID := handler.GetID(ctx)

	address, err := m.LoadRechargeAddressByUser(userID)
	if err != nil {
		return address, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if address == "" {
		_, err = m.TransactionMgr.AllocateTronAddress(userID)
		if err != nil {
			return "", err
		}
	}

	return address, nil
}

// Withdraw allows users to initiate a withdrawal.
func (m *Mall) Withdraw(ctx context.Context, withdrawAddr, value string) error {
	userID := handler.GetID(ctx)

	if withdrawAddr == "" || value == "" {
		return &api.ErrWeb{Code: terrors.ParametersWrong.Int(), Message: terrors.ParametersWrong.String()}
	}

	_, err := utils.ReduceBigInt(value, "0")
	if err != nil {
		return err
	}

	cfg, err := m.GetMallConfigFunc()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return &api.ErrWeb{Code: terrors.ConfigError.Int(), Message: err.Error()}
	}

	node := trxbridge.NewGrpcClient(cfg.TrxHTTPSAddr)
	err = node.Start()
	if err != nil {
		return &api.ErrWeb{Code: terrors.ConfigError.Int(), Message: err.Error()}
	}

	_, err = node.GetAccount(withdrawAddr)
	if err != nil {
		return &api.ErrWeb{Code: terrors.WithdrawAddrError.Int(), Message: err.Error()}
	}

	return m.WithdrawManager.CreateWithdrawOrder(userID, withdrawAddr, value)
}

// GetUserRechargeRecords retrieves user recharge records with pagination.
func (m *Mall) GetUserRechargeRecords(ctx context.Context, limit, page int64) (*types.RechargeResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadRechargeRecordsByUser(userID, limit, page)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

// GetUserWithdrawalRecords retrieves user withdrawal records with pagination.
func (m *Mall) GetUserWithdrawalRecords(ctx context.Context, limit, page int64) (*types.GetWithdrawResponse, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadWithdrawRecordsByUser(userID, limit, page)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return info, nil
}

// GetUserInstanceRecords retrieves user instance records with pagination.
func (m *Mall) GetUserInstanceRecords(ctx context.Context, limit, page int64) (*types.UserInstanceResponse, error) {
	userID := handler.GetID(ctx)
	k, s := m.getAliAccessKeys()
	instanceInfos, err := m.LoadInstancesInfoByUser(userID, limit, page)
	if err != nil {
		log.Errorf("LoadMyInstancesInfo err: %s", err.Error())
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	for _, instanceInfo := range instanceInfos.List {
		var instanceIds []string
		instanceIds = append(instanceIds, instanceInfo.InstanceId)
		log.Infof("InstanceId %s : %s", instanceInfo.RegionId, instanceInfo.InstanceId)
		rsp, err := aliyun.DescribeInstanceStatus(instanceInfo.RegionId, k, s, instanceIds)
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
		instanceExpiredTime, err := aliyun.DescribeInstances(instanceInfo.RegionId, k, s, instanceIds)
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

// GetInstanceDetailsInfo retrieves details of a specific instance.
func (m *Mall) GetInstanceDetailsInfo(ctx context.Context, instanceID string) (*types.InstanceDetails, error) {
	userID := handler.GetID(ctx)

	info, err := m.LoadInstanceInfoByUser(userID, instanceID)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}
	if info.DataDiskString != "" {
		if err := json.Unmarshal([]byte(info.DataDiskString), &info.DataDisk); err != nil {
			return info, nil
		}
	}
	return info, nil
}

// UpdateInstanceName updates the name of a specific instance.
func (m *Mall) UpdateInstanceName(ctx context.Context, instanceID, instanceName string) error {
	userID := handler.GetID(ctx)
	err := m.UpdateVpsInstanceName(instanceID, instanceName, userID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return nil
}

// GetInstanceDefaultInfo retrieves default instance information with pagination.
func (m *Mall) GetInstanceDefaultInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) (*types.InstanceTypeResponse, error) {
	req.Offset = req.Limit * (req.Page - 1)
	instanceInfo, err := m.LoadInstanceDefaultInfo(req)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	usdRate := utils.GetUSDRate()
	for _, info := range instanceInfo.List {
		info.OriginalPrice = info.OriginalPrice / usdRate
		info.Price = info.Price / usdRate
	}
	return instanceInfo, nil
}

// GetInstanceCpuInfo retrieves CPU information for instances.
func (m *Mall) GetInstanceCpuInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) ([]*int32, error) {
	return m.LoadInstanceCPUInfo(req)
}

// GetInstanceMemoryInfo retrieves memory information for instances.
func (m *Mall) GetInstanceMemoryInfo(ctx context.Context, req *types.InstanceTypeFromBaseReq) ([]*float32, error) {
	return m.LoadInstanceMemoryInfo(req)
}

// GetSignCode generates a sign code for a user.
func (m *Mall) GetSignCode(ctx context.Context, userID string) (string, error) {
	return m.UserMgr.GenerateSignCode(userID), nil
}

// Login authenticates a user and generates a JWT token.
func (m *Mall) Login(ctx context.Context, user *types.UserReq) (*types.LoginResponse, error) {
	userID := user.UserId
	log.Debugf("login user:%s", userID)

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

	log.Debugf("login address:%s , %s", address, code)

	p := types.JWTPayload{
		ID:        address,
		LoginType: int64(user.Type),
		Allow:     []auth.Permission{api.RoleUser},
	}
	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	rsp := &types.LoginResponse{}
	rsp.UserId = address
	rsp.Token = string(tk)
	err = m.initUser(address)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

// initUser initializes a user's data if it doesn't exist.
func (m *Mall) initUser(userID string) error {
	exist, err := m.UserExists(userID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}
	if !exist {
		err = m.SaveUserInfo(&types.UserInfo{UserID: userID, Balance: "0"})
		if err != nil {
			return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
		}
	}
	// init recharge address
	addr, err := m.LoadRechargeAddressByUser(userID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}
	if addr == "" {
		_, err = m.TransactionMgr.AllocateTronAddress(userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Logout logs out a user.
func (m *Mall) Logout(ctx context.Context, user *types.UserReq) error {
	userID := handler.GetID(ctx)
	log.Warnf("user id : %s", userID)
	// delete(m.UserMgr.User, user.UserId)
	return nil
}
