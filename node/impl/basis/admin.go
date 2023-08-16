package basis

import (
	"context"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/handler"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
)

func (m *Basis) GetAdminSignCode(ctx context.Context, userId string) (string, error) {
	return m.UserMgr.SetSignCode(userId)
}

func (m *Basis) LoginAdmin(ctx context.Context, user *types.UserReq) (*types.UserResponse, error) {
	userID := user.UserId
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

func (m *Basis) GetWithdrawalRecords(ctx context.Context, limit, offset int64) (*types.WithdrawResponse, error) {
	return m.LoadWithdrawRecords(limit, offset)
}

func (m *Basis) UpdateWithdrawalRecord(ctx context.Context, orderID, withdrawHash string) error {
	userID := handler.GetID(ctx)

	info, err := m.LoadWithdrawRecord(orderID)
	if err != nil {
		return err
	}

	info.WithdrawHash = withdrawHash
	info.Executor = userID

	return m.UpdateWithdrawRecord(info, types.WithdrawDone)
}
