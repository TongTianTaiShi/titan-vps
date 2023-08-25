package user

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"golang.org/x/xerrors"
)

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	user sync.Map // map[string]*types.UserInfoTmp
}

// NewManager creates a new instance of the node manager
func NewManager() (*Manager, error) {
	manager := &Manager{}
	return manager, nil
}

func (m *Manager) GenerateSignCode(key string) string {
	userInfo := &types.UserInfoTmp{}
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	Code := "Vps(" + fmt.Sprintf("%06d", randNew.Intn(1000000)) + ")"
	vI, ok := m.user.Load(key)
	if ok {
		userInfo = vI.(*types.UserInfoTmp)
	} else {
		userInfo.UserLogin.UserId = key
	}
	userInfo.UserLogin.SignCode = Code

	m.user.Store(key, userInfo)
	return Code
}

func (m *Manager) GetSignCode(key string) (string, error) {
	vI, ok := m.user.Load(key)
	userInfo := vI.(*types.UserInfoTmp)
	if ok {
		code := userInfo.UserLogin.SignCode
		if code != "" {
			userInfo.UserLogin.SignCode = ""
			m.user.Store(key, userInfo)
			return code, nil
		}
	}

	return "", xerrors.New("sign code is null")
}
