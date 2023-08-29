package user

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Manager is responsible for managing users.
type Manager struct {
	userCodes sync.Map
}

// NewManager creates a new instance of the user manager.
func NewManager() (*Manager, error) {
	manager := &Manager{}
	return manager, nil
}

// GenerateSignCode generates and stores a sign code for a user.
func (m *Manager) GenerateSignCode(userID string) string {
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := "Vps(" + fmt.Sprintf("%06d", randNew.Intn(1000000)) + ")"

	m.userCodes.Store(userID, code)
	return code
}

// GetSignCode retrieves and removes the sign code for a user.
func (m *Manager) GetSignCode(userID string) (string, error) {
	codeI, ok := m.userCodes.Load(userID)
	if !ok {
		return "", fmt.Errorf("no sign code found for user: %s", userID)
	}

	code, ok := codeI.(string)
	if !ok || code == "" {
		return "", fmt.Errorf("invalid sign code for user: %s", userID)
	}

	m.userCodes.Delete(userID)
	return code, nil
}
