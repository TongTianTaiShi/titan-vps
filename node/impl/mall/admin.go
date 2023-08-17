package mall

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
	"golang.org/x/xerrors"
)

func (m *Mall) GetAdminSignCode(ctx context.Context, userID string) (string, error) {
	exist, err := m.AdminExists(userID)
	if err != nil {
		return "", err
	}

	if !exist {
		return "", xerrors.New("you are not an administrator")
	}

	return m.UserMgr.SetSignCode(userID)
}

func (m *Mall) LoginAdmin(ctx context.Context, user *types.UserReq) (*types.UserResponse, error) {
	userID := user.UserId

	exist, err := m.AdminExists(userID)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, xerrors.New("you are not an administrator")
	}

	code, err := m.UserMgr.GetSignCode(userID)
	if err != nil {
		return nil, err
	}
	signature := user.Signature
	address, err := verifyEthMessage(code, signature)
	if err != nil {
		return nil, err
	}

	p := types.JWTPayload{
		ID:        address,
		LoginType: int64(user.Type),
		Allow:     []auth.Permission{api.RoleAdmin},
	}
	rsp := &types.UserResponse{}
	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return rsp, err
	}
	rsp.UserId = address
	rsp.Token = string(tk)

	return rsp, nil
}

func (m *Mall) GetWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) {
	return m.LoadWithdrawRecords(limit, offset)
}

func (m *Mall) UpdateWithdrawalRecord(ctx context.Context, orderID, withdrawHash string) error {
	userID := handler.GetID(ctx)

	info, err := m.LoadWithdrawRecord(orderID)
	if err != nil {
		return err
	}

	info.WithdrawHash = withdrawHash
	info.Executor = userID

	return m.UpdateWithdrawRecord(info, types.WithdrawDone)
}

func (m *Mall) AddAdminUser(ctx context.Context, userID, nickName string) error {
	return m.SaveAdminInfo(userID, nickName)
}
