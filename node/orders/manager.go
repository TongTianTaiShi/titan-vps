package orders

import (
	"context"
	"sync"

	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/filecoin-project/go-statemachine"
	"github.com/ipfs/go-datastore"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("asset")

const (
	checkLimit = 100
)

// Manager manages asset replicas
type Manager struct {
	stateMachineWait   sync.WaitGroup
	assetStateMachines *statemachine.StateGroup
}

// NewManager returns a new AssetManager instance
func NewManager(ds datastore.Batching, sdb *db.SQLDB) *Manager {
	m := &Manager{}

	// state machine initialization
	m.stateMachineWait.Add(1)
	m.assetStateMachines = statemachine.New(ds, m, OrderInfo{})

	return m
}

// Start initializes and starts the asset state machine and associated tickers
func (m *Manager) Start(ctx context.Context) {
	if err := m.initStateMachines(ctx); err != nil {
		log.Errorf("restartStateMachines err: %s", err.Error())
	}
}

// Terminate stops the asset state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.assetStateMachines.Stop(ctx)
}

func (m *Manager) CreatedOrder() error {
	return nil
}
