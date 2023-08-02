package orders

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinfvm"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/filecoin-project/go-statemachine"
	"github.com/filecoin-project/pubsub"
	"github.com/google/uuid"
	"github.com/ipfs/go-datastore"
	"golang.org/x/xerrors"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("orders")

const (
	checkOrderInterval = 10 * time.Second
	checkLimit         = 100
	orderTimeoutTime   = 5 * time.Minute
)

// Manager manager order
type Manager struct {
	stateMachineWait   sync.WaitGroup
	orderStateMachines *statemachine.StateGroup
	*db.SQLDB

	notify *pubsub.PubSub

	ongoingOrders map[string]string
	orderLock     *sync.Mutex

	usabilityAddrs map[string]string
	usedAddrs      map[string]string
	addrLock       *sync.Mutex
}

// NewManager returns a new manager instance
func NewManager(ds datastore.Batching, sdb *db.SQLDB, pb *pubsub.PubSub) *Manager {
	m := &Manager{
		SQLDB:          sdb,
		notify:         pb,
		ongoingOrders:  make(map[string]string),
		orderLock:      &sync.Mutex{},
		usabilityAddrs: make(map[string]string),
		usedAddrs:      make(map[string]string),
		addrLock:       &sync.Mutex{},
	}

	// state machine initialization
	m.stateMachineWait.Add(1)
	m.orderStateMachines = statemachine.New(ds, m, OrderInfo{})

	return m
}

func (m *Manager) initAddress(as []string) {
	for _, addr := range as {
		m.usabilityAddrs[addr] = ""
	}
}

// Start initializes and starts the order state machine and associated tickers
func (m *Manager) Start(ctx context.Context) {
	// TODO
	m.initAddress([]string{
		"0x5feaAc40B8eB3575794518bC0761cB4A95838ccF",
		"0xddfa8C217a0Fb51a6319e2D863743807d07C9e81",
	})

	if err := m.initStateMachines(ctx); err != nil {
		log.Errorf("restartStateMachines err: %s", err.Error())
	}

	go m.subscribeNodeEvents()
	go m.checkOrderTimeout()
}

func (m *Manager) checkOrderTimeout() {
	ticker := time.NewTicker(checkOrderInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		for addr, hash := range m.usedAddrs {
			info, err := m.LoadOrderRecord(hash)
			if err != nil {
				log.Errorf("checkOrderTimeout LoadOrderRecord %s , %s err:%s", addr, hash, err.Error())
				continue
			}

			log.Debugf("checkout %s , %s ", addr, hash)

			if info.State != int64(Done) && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {
				err = m.orderStateMachines.Send(OrderHash(hash), OrderTimeOut{})
				if err != nil {
					log.Errorf("checkOrderTimeout Send %s , %s err:%s", addr, hash, err.Error())
					continue
				}
			}
		}
	}
}

func (m *Manager) subscribeNodeEvents() {
	subTransfer := m.notify.Sub(types.EventTransfer.String())

	go func() {
		defer m.notify.Unsub(subTransfer)

		for {
			select {
			case u := <-subTransfer:
				tr := u.(*filecoinfvm.AbiTransfer)
				log.Debugf("to hex", tr.To.Hex())
				if hash, exist := m.usedAddrs[tr.To.Hex()]; exist {
					err := m.orderStateMachines.Send(OrderHash(hash), PaymentSucceed{})
					if err != nil {
						log.Errorf("subscribeNodeEvents Send %s err:%s", hash, err.Error())
						continue
					}
				}
			}
		}
	}()
}

// Terminate stops the order state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.orderStateMachines.Stop(ctx)
}

// CancelOrder cancel vps order
func (m *Manager) CancelOrder(orderID string) error {
	return m.orderStateMachines.Send(OrderHash(orderID), OrderCancel{})
}

// CreatedOrder create vps order
func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	err := m.addOrder(req.From, orderID)
	if err != nil {
		return err
	}

	address, err := m.allocatePayeeAddress(orderID)
	if err != nil {
		return err
	}

	req.To = address
	req.OrderID = orderID

	err = m.SaveOrderRecord(req)
	if err != nil {
		return err
	}

	// create order task
	return m.orderStateMachines.Send(OrderHash(orderID), WaitingPaymentSent{})
}

func (m *Manager) addOrder(userID, orderID string) error {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	if _, exist := m.ongoingOrders[userID]; exist {
		return xerrors.New("user have order")
	}

	m.ongoingOrders[userID] = orderID

	return nil
}

func (m *Manager) removeOrder(userID string) {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	delete(m.ongoingOrders, userID)
}

func (m *Manager) allocatePayeeAddress(orderID string) (string, error) {
	m.addrLock.Lock()
	defer m.addrLock.Unlock()

	if len(m.usabilityAddrs) > 0 {
		for addr := range m.usabilityAddrs {
			m.usedAddrs[addr] = orderID
			delete(m.usabilityAddrs, addr)
			return addr, nil
		}
	}

	return "", xerrors.New("not found address")
}

func (m *Manager) revertPayeeAddress(addr string) {
	m.addrLock.Lock()
	defer m.addrLock.Unlock()

	delete(m.usedAddrs, addr)
	m.usabilityAddrs[addr] = ""
}

func (m *Manager) recoverOutstandingOrders(info OrderInfo) {
	m.addrLock.Lock()
	defer m.addrLock.Unlock()

	m.usedAddrs[info.To] = info.OrderID.String()
	delete(m.usabilityAddrs, info.To)

	m.addOrder(info.From, info.OrderID.String())
}
