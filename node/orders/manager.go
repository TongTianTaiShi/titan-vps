package orders

import (
	"context"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/LMF709268224/titan-vps/node/vps"
	"github.com/filecoin-project/go-statemachine"
	"github.com/filecoin-project/pubsub"
	"github.com/ipfs/go-datastore"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("orders")

const (
	checkOrderInterval = 10 * time.Second
	orderTimeoutMinute = 10
	orderTimeoutTime   = orderTimeoutMinute * time.Minute
)

// OrderManager manages order processing.
type Manager struct {
	stateMachineWait   sync.WaitGroup
	orderStateMachines *statemachine.StateGroup
	*db.SQLDB

	notification *pubsub.PubSub

	activeOrders sync.Map // map[string]*types.OrderRecord

	cfg    config.MallCfg
	txMgr  *transaction.Manager
	vpsMgr *vps.Manager
}

// NewOrderManager creates a new order manager instance.
func NewManager(ds datastore.Batching, sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, fm *transaction.Manager, vm *vps.Manager) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		SQLDB:        sdb,
		notification: pb,
		cfg:          cfg,
		txMgr:        fm,
		vpsMgr:       vm,
	}

	// state machine initialization
	m.stateMachineWait.Add(1)
	m.orderStateMachines = statemachine.New(ds, m, OrderInfo{})

	return m, nil
}

// Start initializes and starts the order state machine and related processes.
func (m *Manager) Start(ctx context.Context) {
	if err := m.initStateMachines(ctx); err != nil {
		log.Errorf("restartStateMachines err: %s", err.Error())
	}

	// go m.subscribeEvents()
	go m.checkOrdersTimeout()
}

func (m *Manager) checkOrdersTimeout() {
	ticker := time.NewTicker(checkOrderInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		m.activeOrders.Range(func(key, value interface{}) bool {
			orderRecord := value.(*types.OrderRecord)
			orderID := orderRecord.OrderID

			info, err := m.LoadOrderRecord(orderID, orderTimeoutMinute)
			if err != nil {
				log.Errorf("checkOrderTimeout LoadOrderRecord , %s err:%s", orderID, err.Error())
				return true
			}

			log.Debugf("checkout  %s ", orderID)

			if info.State.Int() != Done.Int() && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {

				err = m.orderStateMachines.Send(OrderHash(orderID), OrderTimeOut{})
				if err != nil {
					log.Errorf("checkOrderTimeout Send  %s err:%s", orderID, err.Error())
					return true
				}
			}
			return true
		})

	}
}

// Terminate stops the order state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.orderStateMachines.Stop(ctx)
}

// CancelOrder cancels a VPS order.
func (m *Manager) CancelOrder(orderID, userID string) error {
	order, err := m.LoadOrderRecord(orderID, orderTimeoutMinute)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if order.UserID != userID {
		return &api.ErrWeb{Code: terrors.UserMismatch.Int(), Message: terrors.UserMismatch.String()}
	}

	err = m.orderStateMachines.Send(OrderHash(orderID), OrderCancel{})
	if err != nil {
		return &api.ErrWeb{Code: terrors.StateMachinesError.Int(), Message: err.Error()}
	}

	return nil
}

// PaymentCompleted marks a VPS order as completed.
func (m *Manager) PaymentCompleted(orderID, userID string) error {
	order, err := m.LoadOrderRecord(orderID, orderTimeoutMinute)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	if order.UserID != userID {
		return &api.ErrWeb{Code: terrors.UserMismatch.Int(), Message: terrors.UserMismatch.String()}
	}

	err = m.orderStateMachines.Send(OrderHash(orderID), PaymentResult{})
	if err != nil {
		return &api.ErrWeb{Code: terrors.StateMachinesError.Int(), Message: err.Error()}
	}

	return nil
}

// CreatedOrder creates a new VPS order.
func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()

	m.addOrder(req)

	// create order task
	err := m.orderStateMachines.Send(OrderHash(req.OrderID), CreateOrder{orderInfoFrom(req)})
	if err != nil {
		return &api.ErrWeb{Code: terrors.StateMachinesError.Int(), Message: err.Error()}
	}

	return nil
}

func (m *Manager) addOrder(req *types.OrderRecord) {
	m.activeOrders.Store(req.OrderID, req)
}

func (m *Manager) removeOrder(orderID string) {
	m.activeOrders.Delete(orderID)
}

func (m *Manager) getHeight() int64 {
	var msg filecoinbridge.TipSet
	err := filecoinbridge.ChainHead(&msg, m.cfg.LotusHTTPSAddr)
	if err != nil {
		log.Errorf("ChainHead err:%s", err.Error())
		return 0
	}

	return msg.Height
}

// GetOrderTimeoutDurationMinutes returns the duration in minutes after which an order times out.
func (m *Manager) GetOrderTimeoutDurationMinutes() int {
	return orderTimeoutMinute
}
