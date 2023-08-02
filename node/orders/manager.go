package orders

import (
	"context"
	"strings"
	"sync"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/filecoin-project/go-statemachine"
	"github.com/google/uuid"
	"github.com/ipfs/go-datastore"
	"golang.org/x/xerrors"

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
	*db.SQLDB

	userOrderMap  map[string]string
	userOrderLock *sync.Mutex
}

// NewManager returns a new AssetManager instance
func NewManager(ds datastore.Batching, sdb *db.SQLDB) *Manager {
	m := &Manager{
		SQLDB:         sdb,
		userOrderMap:  make(map[string]string),
		userOrderLock: &sync.Mutex{},
	}

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

func (m *Manager) CancelOrder(hash string) error {
	return m.assetStateMachines.Send(OrderHash(hash), OrderCancel{})
}

func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()

	hash := uuid.NewString()
	hash = strings.Replace(hash, "-", "", -1)

	err := m.addOrderToUser(req.From, hash)
	if err != nil {
		return err
	}

	req.Hash = hash

	err = m.SaveOrderRecord(req)
	if err != nil {
		return err
	}

	// create asset task
	return m.assetStateMachines.Send(OrderHash(hash), WaitingPaymentSent{})
}

func (m *Manager) addOrderToUser(user, order string) error {
	m.userOrderLock.Lock()
	defer m.userOrderLock.Unlock()

	if _, exist := m.userOrderMap[user]; exist {
		return xerrors.New("user have order")
	}

	m.userOrderMap[user] = order

	return nil
}

func (m *Manager) deleteOrderFromUser(user string) {
	m.userOrderLock.Lock()
	defer m.userOrderLock.Unlock()

	delete(m.userOrderMap, user)
}
