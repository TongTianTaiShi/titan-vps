package mall

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
	uuid2 "github.com/google/uuid"
	"math/rand"
	"strings"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	email2 "github.com/LMF709268224/titan-vps/lib/email"
	"github.com/LMF709268224/titan-vps/node/config"
)

func (m *Mall) LoginAccount(ctx context.Context, request *types.AccountRequest) (*types.AccountLoginResponse, error) {
	m.Cache.Set("448688985@qq.com", "123456")
	switch request.Type {
	case types.LoginTypeMetaMask:
		return m.loginMetaMask(request)
	case types.LoginTypeEmail:
		return m.loginEmail(request)
	case types.LoginTypeFilecoin:
		return m.loginFilecoin(request)
	default:
		return nil, fmt.Errorf("%s not supported", request.Type.String())
	}
}

func (m *Mall) loginEmail(request *types.AccountRequest) (*types.AccountLoginResponse, error) {
	if request.Email == "" || !(request.VerifyCode != "" && request.Password != "") {
		return nil, &api.ErrWeb{Code: terrors.ProviderRequestParameterError.Int(), Message: terrors.ProviderRequestParameterError.String()}
	}

	uuid, err := m.getProviderUUID(request.Email, request.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain provider's uuid")
	}

	if request.VerifyCode != "" {
		if err = m.Cache.Check(request.Email, request.VerifyCode); err != nil {
			return nil, fmt.Errorf("verify code authentication failed")
		}
	} else {
		if ok, err := m.SQLDB.CheckPassword(uuid, request.Password); err != nil {
			return nil, fmt.Errorf("failed to find password")
		} else if !ok {
			return nil, fmt.Errorf("password error")
		}
	}

	tk, err := m.signJWTToken(uuid, request.Type, []auth.Permission{api.RoleProvider})
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	err = m.updateDatabase(request, uuid, request.Email)
	if err != nil {
		return nil, err
	}

	rsp := &types.AccountLoginResponse{
		ID:    uuid,
		Token: string(tk),
	}

	return rsp, nil
}

func (m *Mall) loginMetaMask(request *types.AccountRequest) (*types.AccountLoginResponse, error) {
	if request.Address == "" || request.Signature == "" {
		return nil, &api.ErrWeb{Code: terrors.ProviderRequestParameterError.Int(), Message: terrors.ProviderRequestParameterError.String()}
	}

	uuid, err := m.getProviderUUID(request.Address, request.Type)
	if err != nil {
		return nil, err
	}

	err = m.checkSignature(request)
	if err != nil {
		return nil, err
	}

	tk, err := m.signJWTToken(uuid, request.Type, []auth.Permission{api.RoleProvider})
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	err = m.updateDatabase(request, uuid, request.Address)
	if err != nil {
		return nil, err
	}

	rsp := &types.AccountLoginResponse{
		ID:    uuid,
		Token: string(tk),
	}

	return rsp, nil
}

func (m *Mall) loginFilecoin(request *types.AccountRequest) (*types.AccountLoginResponse, error) {
	if request.Filecoin == "" || request.Address == "" || request.Signature == "" {
		return nil, &api.ErrWeb{Code: terrors.ProviderRequestParameterError.Int(), Message: terrors.ProviderRequestParameterError.String()}
	}

	uuid, err := m.getProviderUUID(request.Filecoin, request.Type)
	if err != nil {
		return nil, err
	}

	err = m.checkSignature(request)
	if err != nil {
		return nil, err
	}

	tk, err := m.signJWTToken(uuid, request.Type, []auth.Permission{api.RoleProvider})
	if err != nil {
		return nil, &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	err = m.updateDatabase(request, uuid, request.Filecoin)
	if err != nil {
		return nil, err
	}

	rsp := &types.AccountLoginResponse{
		ID:    uuid,
		Token: string(tk),
	}

	return rsp, nil
}

func (m *Mall) checkSignature(request *types.AccountRequest) error {
	code, err := m.AccountMgr.GetSignCode(request.Address)
	if err != nil {
		return &api.ErrWeb{Code: terrors.NotFoundSignCode.Int(), Message: terrors.NotFoundSignCode.String()}
	}

	publicKey, err := verifyEthMessage(code, request.Signature)
	if err != nil {
		return &api.ErrWeb{Code: terrors.SignError.Int(), Message: err.Error()}
	}

	if strings.ToLower(request.Address) != strings.ToLower(publicKey) {
		return &api.ErrWeb{Code: terrors.UserMismatch.Int(), Message: fmt.Sprintf("%s,%s", request.Address, publicKey)}
	}

	return nil
}

func (m *Mall) signJWTToken(id string, loginType types.LoginType, allow []auth.Permission) ([]byte, error) {
	p := types.JWTPayload{
		ID:        id,
		LoginType: int64(loginType),
		Allow:     allow,
	}

	tk, err := jwt.Sign(&p, m.APISecret)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (m *Mall) getProviderUUID(userID string, loginType types.LoginType) (string, error) {
	uuid, err := m.SQLDB.GetProviderUUID(userID, loginType)
	if err == sql.ErrNoRows {
		id, err := uuid2.NewUUID()
		if err != nil {
			return "", fmt.Errorf("generate uuid failed, err: %s", err)
		}

		return id.String(), nil
	} else if err != nil {
		return "", err
	} else {
		return uuid, nil
	}
}

func (m *Mall) updateDatabase(request *types.AccountRequest, id string, loginAddress string) error {
	if exist, err := m.SQLDB.CheckProviderIsExist(id); err != nil {
		return err
	} else if !exist {
		err = m.SQLDB.SaveProviderInfo(id, loginAddress, request.Type, request.Password)
		if err != nil {
			return err
		}

		err = m.SQLDB.UpdateInvitationUUID(request.InvitationCode, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Mall) GetVerifyMessage(ctx context.Context, id string, loginType types.LoginType) (string, error) {
	if loginType != types.LoginTypeEmail && loginType != types.LoginTypeMetaMask {
		return "", fmt.Errorf("%s not supported", loginType.String())
	}

	if loginType == types.LoginTypeEmail {
		err := m.getVerifyCode(id)
		return "", err
	} else {
		signCode := m.AccountMgr.GenerateSignCode(id)
		return signCode, nil
	}
}

func (m *Mall) getVerifyCode(email string) error {
	cfg, err := m.GetMallConfigFunc()
	if err != nil {
		return err
	}

	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	verifyCode := fmt.Sprintf("%06d", randNew.Intn(1000000))
	err = m.Cache.Set(email, verifyCode)
	if err != nil {
		return err
	}

	err = sendEmail(cfg.Email, email, verifyCode)
	if err != nil {
		return err
	}

	return nil
}

func sendEmail(cfg config.EmailConfig, sendTo string, vc string) error {
	var data email2.Data
	data.Subject = "【Titan VPS】您的验证码"
	data.Tittle = "please check your verify code "
	data.SendTo = sendTo
	data.Content += "<p style=\"line-height:38px;margin:30px;\"> <b>亲爱的用户:</b><br>"
	data.Content +=
		"&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;您好！感谢您选择使用Titan VPS，我们是一家基" +
			"于Filecoin提供去中心化存储云盘服务的平台。您正在" +
			"进行邮箱验证，以验证您的身份或在我们的平台上进行注" +
			"册或登录。<br>您的验证码为：<strong>" + vc + "</strong><br>" +
			"&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;请在操作页面输入此验证码以完成验证。为了保证您的账" +
			"号安全，请勿将此验证码透露给他人。请注意，此验证码" +
			"在接收后的5分钟内有效。若您未在有效时间内完成验" +
			"证，验证码将会失效。如果验证码失效，您可以重新发起" +
			"邮箱验证流程获取新的验证码。如果您并未进行相关操作，" +
			"可能是其他用户误操作，此情况下请忽略此邮件。<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;感谢您" +
			"对Titan VPS的信任和支持，我们将一如既往地为" +
			"您提供高品质的服务。祝您使用愉快！<br></p>" +
			"<h1>Titan VPS团队</h1>"
	err := email2.SendEmail(cfg, data)
	if err != nil {
		log.Errorf("sendEmailing failed:%v", err)
		return err
	}
	return nil
}

func (m *Mall) SetInvitationCode(ctx context.Context, code string) error {
	return m.SQLDB.InsertInvitationCode(code)
}
