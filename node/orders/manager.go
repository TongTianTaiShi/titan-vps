package orders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/aliyun"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
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

	cfg  config.BasisCfg
	tMgr *transaction.Manager
}

// NewManager returns a new manager instance
func NewManager(ds datastore.Batching, sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetBasisConfigFunc, fm *transaction.Manager) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &Manager{
		SQLDB:          sdb,
		notify:         pb,
		ongoingOrders:  make(map[string]string),
		orderLock:      &sync.Mutex{},
		usabilityAddrs: make(map[string]string),
		usedAddrs:      make(map[string]string),
		addrLock:       &sync.Mutex{},
		cfg:            cfg,
		tMgr:           fm,
	}

	// state machine initialization
	m.stateMachineWait.Add(1)
	m.orderStateMachines = statemachine.New(ds, m, OrderInfo{})

	return m, nil
}

func (m *Manager) initPaymentAddress(as []string) {
	m.addrLock.Lock()
	defer m.addrLock.Unlock()

	for _, addr := range as {
		m.usabilityAddrs[addr] = ""
	}
}

// Start initializes and starts the order state machine and associated tickers
func (m *Manager) Start(ctx context.Context) {
	m.initPaymentAddress(m.cfg.PaymentAddress)

	if err := m.initStateMachines(ctx); err != nil {
		log.Errorf("restartStateMachines err: %s", err.Error())
	}

	go m.subscribeEvents()
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

			if info.State != Done.Int() && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {

				height := m.tMgr.GetHeight()

				err = m.orderStateMachines.Send(OrderHash(hash), OrderTimeOut{Height: height})
				if err != nil {
					log.Errorf("checkOrderTimeout Send %s , %s err:%s", addr, hash, err.Error())
					continue
				}
			}
		}
	}
}

func (m *Manager) subscribeEvents() {
	subTransfer := m.notify.Sub(types.EventTransferWatch.String())
	defer m.notify.Unsub(subTransfer)

	for {
		select {
		case u := <-subTransfer:
			tr := u.(*types.FvmTransferWatch)

			if hash, exist := m.usedAddrs[tr.To]; exist {
				err := m.orderStateMachines.Send(OrderHash(hash), PaymentResult{
					&PaymentInfo{
						ID:    tr.ID,
						From:  tr.From,
						To:    tr.To,
						Value: tr.Value,
					},
				})
				if err != nil {
					log.Errorf("subscribeNodeEvents Send %s err:%s", hash, err.Error())
					continue
				}
			}
		}
	}
}

// Terminate stops the order state machine
func (m *Manager) Terminate(ctx context.Context) error {
	return m.orderStateMachines.Stop(ctx)
}

// CancelOrder cancel vps order
func (m *Manager) CancelOrder(orderID string) error {
	height := m.tMgr.GetHeight()

	return m.orderStateMachines.Send(OrderHash(orderID), OrderCancel{Height: height})
}

// CreatedOrder create vps order
func (m *Manager) CreatedOrder(req *types.OrderRecord) error {
	m.stateMachineWait.Wait()

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	err := m.addOrder(req.User, orderID)
	if err != nil {
		return err
	}

	address, err := m.allocatePayeeAddress(orderID)
	if err != nil {
		return err
	}

	req.To = address
	req.OrderID = orderID
	req.CreatedHeight = m.tMgr.GetHeight()

	// create order task
	return m.orderStateMachines.Send(OrderHash(orderID), CreateOrder{orderInfoFrom(req)})
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

	m.addOrder(info.User, info.OrderID.String())
}

func (m *Manager) createAliyunInstance(vpsInfo *types.CreateInstanceReq) (*types.CreateInstanceResponse, error) {
	k := m.cfg.AliyunAccessKeyID
	s := m.cfg.AliyunAccessKeySecret

	priceUnit := vpsInfo.PeriodUnit
	period := vpsInfo.Period
	regionID := vpsInfo.RegionId
	instanceType := vpsInfo.InstanceType
	imageID := vpsInfo.ImageId
	if priceUnit == "Year" {
		priceUnit = "Month"
		period = period * 12
	}

	var securityGroupID string

	group, err := aliyun.DescribeSecurityGroups(regionID, k, s)
	if err == nil && len(group) > 0 {
		securityGroupID = group[0]
	} else {
		securityGroupID, err = aliyun.CreateSecurityGroup(regionID, k, s)
		if err != nil {
			log.Errorf("CreateSecurityGroup err: %s", err.Error())
			return nil, xerrors.New(*err.Data)
		}
	}

	log.Debugln("securityGroupID:", securityGroupID, " , DryRun:", vpsInfo.DryRun)

	result, err := aliyun.CreateInstance(regionID, k, s, instanceType, imageID, securityGroupID, priceUnit, period, vpsInfo.DryRun)
	if err != nil {
		log.Errorf("CreateInstance err: %s", err.Error())
		return nil, xerrors.New(*err.Data)
	}

	address, err := aliyun.AllocatePublicIPAddress(regionID, k, s, result.InstanceID)
	if err != nil {
		log.Errorf("AllocatePublicIpAddress err: %s", err.Error())
	} else {
		result.PublicIpAddress = address
	}

	err = aliyun.AuthorizeSecurityGroup(regionID, k, s, securityGroupID)
	if err != nil {
		log.Errorf("AuthorizeSecurityGroup err: %s", err.Error())
	}
	randNew := rand.New(rand.NewSource(time.Now().UnixNano()))
	keyPairName := "KeyPair" + fmt.Sprintf("%06d", randNew.Intn(1000000))
	keyInfo, err := aliyun.CreateKeyPair(regionID, k, s, keyPairName)
	if err != nil {
		log.Errorf("CreateKeyPair err: %s", err.Error())
	} else {
		result.PrivateKey = keyInfo.PrivateKeyBody
	}
	var instanceIds []string
	instanceIds = append(instanceIds, result.InstanceID)
	_, err = aliyun.AttachKeyPair(regionID, k, s, keyPairName, instanceIds)
	if err != nil {
		log.Errorf("AttachKeyPair err: %s", err.Error())
	}
	go func() {
		time.Sleep(1 * time.Minute)

		err := aliyun.StartInstance(regionID, k, s, result.InstanceID)
		log.Infoln("StartInstance err:", err)
	}()

	return result, nil
}
