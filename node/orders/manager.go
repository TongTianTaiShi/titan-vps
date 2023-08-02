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

var log = logging.Logger("asset")

const (
	checkOrderInterval = 10 * time.Second
	checkLimit         = 100
	orderTimeoutTime   = 5 * time.Minute
)

// Manager manages asset replicas
type Manager struct {
	stateMachineWait   sync.WaitGroup
	assetStateMachines *statemachine.StateGroup
	*db.SQLDB

	notify *pubsub.PubSub

	userOrderMap  map[string]string
	userOrderLock *sync.Mutex

	usabilityAddrs map[string]string
	usedAddrs      map[string]string
	addrLock       *sync.Mutex
}

// NewManager returns a new AssetManager instance
func NewManager(ds datastore.Batching, sdb *db.SQLDB, pb *pubsub.PubSub) *Manager {
	m := &Manager{
		SQLDB:          sdb,
		notify:         pb,
		userOrderMap:   make(map[string]string),
		userOrderLock:  &sync.Mutex{},
		usabilityAddrs: make(map[string]string),
		usedAddrs:      make(map[string]string),
		addrLock:       &sync.Mutex{},
	}

	m.usabilityAddrs["0x5feaAc40B8eB3575794518bC0761cB4A95838ccF"] = ""
	m.usabilityAddrs["0xddfa8C217a0Fb51a6319e2D863743807d07C9e81"] = ""

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
				log.Errorf("LoadOrderRecord %s , %s err:%s", addr, hash, err.Error())
				continue
			}

			log.Debugf("checkout %s , %s ", addr, hash)

			if info.State != int64(Done) && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {
				m.assetStateMachines.Send(OrderHash(hash), OrderTimeOut{})
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
					m.assetStateMachines.Send(OrderHash(hash), PaymentSucceed{})
				}
			}
		}
	}()
}

// Terminate stops the asset state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.assetStateMachines.Stop(ctx)
}

func (m *Manager) CancelOrder(orderID string) error {
	return m.assetStateMachines.Send(OrderHash(orderID), OrderCancel{})
}

func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	err := m.addOrderToUser(req.From, orderID)
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

	// create asset task
	return m.assetStateMachines.Send(OrderHash(orderID), WaitingPaymentSent{})
}

func (m *Manager) addOrderToUser(user, orderID string) error {
	m.userOrderLock.Lock()
	defer m.userOrderLock.Unlock()

	if _, exist := m.userOrderMap[user]; exist {
		return xerrors.New("user have order")
	}

	m.userOrderMap[user] = orderID

	return nil
}

func (m *Manager) deleteOrderFromUser(user string) {
	m.userOrderLock.Lock()
	defer m.userOrderLock.Unlock()

	delete(m.userOrderMap, user)
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

	m.addOrderToUser(info.From, info.OrderID.String())
}
