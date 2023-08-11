package exchange

import (
	"strings"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/filecoinbridge"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/filecoin-project/pubsub"
	"github.com/google/uuid"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("exchange")

const (
	checkOrderInterval = 10 * time.Second
	orderTimeoutTime   = 10 * time.Minute
)

// RechargeManager manager recharge order
type RechargeManager struct {
	*db.SQLDB
	cfg    config.BasisCfg
	notify *pubsub.PubSub

	ongoingOrders map[string]*types.RechargeRecord
	orderLock     *sync.Mutex
	tMgr          *transaction.Manager
}

// NewRechargeManager returns a new manager instance
func NewRechargeManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetBasisConfigFunc, fm *transaction.Manager) (*RechargeManager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &RechargeManager{
		SQLDB:  sdb,
		notify: pb,
		cfg:    cfg,

		tMgr: fm,

		ongoingOrders: make(map[string]*types.RechargeRecord),
		orderLock:     &sync.Mutex{},
	}

	m.initOngeingOrders()

	go m.checkOrdersTimeout()
	go m.subscribeEvents()

	return m, nil
}

func (m *RechargeManager) subscribeEvents() {
	subTransfer := m.notify.Sub(types.EventTronTransferWatch.String())
	defer m.notify.Unsub(subTransfer)

	for {
		select {
		case u := <-subTransfer:
			tr := u.(*types.TronTransferWatch)

			m.handleTronTransfer(tr)
		}
	}
}

func (m *RechargeManager) handleTronTransfer(tr *types.TronTransferWatch) {
	orderID := tr.OrderID

	info, err := m.LoadRechargeRecord(orderID)
	if err != nil {
		log.Errorf("handleTronTransfer LoadOrderRecord %s , %s err:%s", tr.TxHash, orderID, err.Error())
		return
	}

	if info.State != types.ExchangeCreated {
		log.Errorf("handleTronTransfer Invalid order status %d , %s", info.State, orderID)
		return
	}

	info.Value = tr.Value
	info.TxHash = tr.TxHash
	info.From = tr.From
	info.DoneHeight = tr.Height

	state := types.ExchangeFail
	info.Msg = tr.State.String()

	if tr.State == core.Transaction_Result_SUCCESS {
		state = types.ExchangeDone

		client := filecoinbridge.NewGrpcClient(m.cfg.LotusHTTPSAddr, m.cfg.TitanContractorAddr)
		hash, err := client.Mint(m.cfg.PrivateKeyStr, info.RechargeAddr, tr.Value)
		if err != nil {
			info.Msg = err.Error()
		} else {
			info.RechargeHash = hash
		}
	}

	err = m.changeOrderState(state, info)
	if err != nil {
		log.Errorf("handleTronTransfer changeOrderState %s err:%s", orderID, err.Error())
		return
	}
}

func (m *RechargeManager) initOngeingOrders() {
	records, err := m.LoadRechargeRecords(types.ExchangeCreated)
	if err != nil {
		log.Errorln("LoadRechargeRecords err:", err.Error())
		return
	}

	for _, info := range records {
		m.tMgr.RecoverOutstandingTronOrders(info.To, info.OrderID)
		m.addOrder(info)
	}
}

func (m *RechargeManager) checkOrdersTimeout() {
	ticker := time.NewTicker(checkOrderInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		for _, orderRecord := range m.ongoingOrders {
			orderID := orderRecord.OrderID
			addr := orderRecord.To

			info, err := m.LoadRechargeRecord(orderID)
			if err != nil {
				log.Errorf("checkOrderTimeout LoadOrderRecord %s , %s err:%s", addr, orderID, err.Error())
				continue
			}

			log.Debugf("checkout %s , %s ", addr, orderID)

			if info.State == types.ExchangeCreated && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {
				err := m.changeOrderState(types.ExchangeTimeout, info)
				if err != nil {
					log.Errorf("checkOrderTimeout UpdateRechargeRecord %s , %s err:%s", addr, orderID, err.Error())
					continue
				}
			}
		}
	}
}

// func (m *RechargeManager) getOrderIDByToAddress(to string) (string, bool) {
// 	for _, orderRecord := range m.ongoingOrders {
// 		if orderRecord.To == to {
// 			return orderRecord.OrderID, true
// 		}
// 	}

// 	return "", false
// }

// CancelRechargeOrder cancel the order
func (m *RechargeManager) CancelRechargeOrder(orderID string) error {
	info, err := m.LoadRechargeRecord(orderID)
	if err != nil {
		return err
	}
	if info.State != types.ExchangeCreated {
		return xerrors.Errorf("Invalid order status %d", info.State)
	}

	return m.changeOrderState(types.ExchangeCancel, info)
}

func (m *RechargeManager) changeOrderState(state types.ExchangeState, info *types.RechargeRecord) error {
	info.DoneHeight = getTronHeight(m.cfg.TrxHTTPSAddr)

	err := m.UpdateRechargeRecord(info, state)
	if err != nil {
		return err
	}

	m.removeOrder(info.User)
	m.tMgr.RevertTronAddress(info.To)

	return nil
}

// CreateRechargeOrder create a recharge order
func (m *RechargeManager) CreateRechargeOrder(userAddr, rechargeAddr string) (addr string, err error) {
	defer func() {
		if err != nil {
			m.removeOrder(userAddr)
			m.tMgr.RevertTronAddress(addr)
		}
	}()

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	info := &types.RechargeRecord{
		OrderID:       orderID,
		User:          userAddr,
		RechargeAddr:  rechargeAddr,
		CreatedHeight: getTronHeight(m.cfg.TrxHTTPSAddr),
	}

	err = m.addOrder(info)
	if err != nil {
		return "", err
	}

	addr, err = m.tMgr.AllocateTronAddress(orderID)
	if err != nil {
		return "", err
	}

	info.To = addr

	err = m.SaveRechargeInfo(info)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func (m *RechargeManager) addOrder(info *types.RechargeRecord) error {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	if _, exist := m.ongoingOrders[info.User]; exist {
		return xerrors.New("user have order")
	}

	m.ongoingOrders[info.User] = info

	return nil
}

func (m *RechargeManager) removeOrder(userID string) {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	delete(m.ongoingOrders, userID)
}
