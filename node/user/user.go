package user

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/xerrors"
)

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	userCodes sync.Map
}

// NewManager creates a new instance of the node manager
func NewManager() (*Manager, error) {
	manager := &Manager{}
	return manager, nil
}

func (m *Manager) GenerateSignCode(userID string) string {
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := "Vps(" + fmt.Sprintf("%06d", randNew.Intn(1000000)) + ")"

	m.userCodes.Store(userID, code)
	return code
}

func (m *Manager) GetSignCode(userID string) (string, error) {
	vI, ok := m.userCodes.Load(userID)
	if ok {
		code := vI.(string)
		if code != "" {
			m.userCodes.Delete(userID)

			return code, nil
		}
	}

	return "", xerrors.New("sign code is null")
}
