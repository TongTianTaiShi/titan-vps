package user

import "github.com/LMF709268224/titan-vps/api/types"

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	User map[string]*types.UserInfoTmp
}

// NewManager creates a new instance of the node manager
func NewManager() (*Manager, error) {
	manager := &Manager{}

	return manager, nil
}
