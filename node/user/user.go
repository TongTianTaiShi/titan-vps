package user

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
	"golang.org/x/xerrors"
	"math/rand"
	"time"
)

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	User map[string]*types.UserInfoTmp
}

// NewManager creates a new instance of the node manager
func NewManager() (*Manager, error) {
	mgr := make(map[string]*types.UserInfoTmp)
	manager := &Manager{
		User: mgr,
	}
	return manager, nil
}
func (m *Manager) SetSignCode(key string) (string, error) {
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	Code := "Vps(" + fmt.Sprintf("%06d", randNew.Intn(1000000)) + ")"
	v, ok := m.User[key]
	if ok {
		v.UserLogin.SignCode = Code
	} else {
		m.User[key] = &types.UserInfoTmp{}
		m.User[key].UserLogin.SignCode = Code
		m.User[key].UserLogin.UserId = key
	}
	return Code, nil
}

func (m *Manager) GetSignCode(key string) (string, error) {
	v, ok := m.User[key]
	fmt.Println(ok)
	fmt.Println(v)
	if ok && v.UserLogin.SignCode != "" {
		code := v.UserLogin.SignCode
		v.UserLogin.SignCode = ""
		return code, nil
	}
	return "", xerrors.New("sign code is null")
}
