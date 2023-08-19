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
	user     map[string]*types.UserInfoTmp
	userLock *sync.Mutex
}

// NewManager creates a new instance of the node manager
func NewManager() (*Manager, error) {
	mgr := make(map[string]*types.UserInfoTmp)
	manager := &Manager{
		user:     mgr,
		userLock: &sync.Mutex{},
	}
	return manager, nil
}

func (m *Manager) GenerateSignCode(key string) string {
	m.userLock.Lock()
	defer m.userLock.Unlock()

	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	Code := "Vps(" + fmt.Sprintf("%06d", randNew.Intn(1000000)) + ")"
	v, ok := m.user[key]
	if ok {
		v.UserLogin.SignCode = Code
	} else {
		m.user[key] = &types.UserInfoTmp{}
		m.user[key].UserLogin.SignCode = Code
		m.user[key].UserLogin.UserId = key
	}
	return Code
}

func (m *Manager) GetSignCode(key string) (string, error) {
	v, ok := m.user[key]
	fmt.Println(ok)
	fmt.Println(v)
	if ok && v.UserLogin.SignCode != "" {
		code := v.UserLogin.SignCode
		v.UserLogin.SignCode = ""
		return code, nil
	}
	return "", xerrors.New("sign code is null")
}
