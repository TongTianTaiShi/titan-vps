package exchange

import (
	"strings"
	"sync"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/LMF709268224/titan-vps/node/utils"
	"github.com/filecoin-project/pubsub"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

// WithdrawManager manager withdraw order
type WithdrawManager struct {
	*db.SQLDB
	cfg    config.BasisCfg
	notify *pubsub.PubSub

	ongoingOrders map[string]*types.WithdrawRecord
	orderLock     *sync.Mutex
	tMgr          *transaction.Manager
}

// NewWithdrawManager returns a new manager instance
func NewWithdrawManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetBasisConfigFunc, fm *transaction.Manager) (*WithdrawManager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &WithdrawManager{
		SQLDB:  sdb,
		notify: pb,
		cfg:    cfg,

		tMgr: fm,

		ongoingOrders: make(map[string]*types.WithdrawRecord),
		orderLock:     &sync.Mutex{},
	}

	m.initOngeingOrders()

	// go m.checkOrdersTimeout()
	go m.subscribeEvents()

	return m, nil
}

func (m *WithdrawManager) subscribeEvents() {
	subTransfer := m.notify.Sub(types.EventFvmTransferWatch.String())
	defer m.notify.Unsub(subTransfer)

	for {
		select {
		case u := <-subTransfer:
			tr := u.(*types.FvmTransferWatch)

			log.Debugf("subscribeEvents tr %s ", tr.To)
			if orderID, exist := m.getOrderIDByToAddress(tr.To); exist {
				log.Debugf("getOrderIDByToAddress orderID %s ", orderID)
				m.handleFvmTransfer(orderID, tr)
			}
		}
	}
}

func (m *WithdrawManager) handleFvmTransfer(orderID string, tr *types.FvmTransferWatch) {
	info, err := m.LoadWithdrawRecord(orderID)
	if err != nil {
		log.Errorf("handleFvmTransfer LoadOrderRecord %s , %s err:%s", tr.To, orderID, err.Error())
		return
	}

	// if info.State != types.ExchangeCreated {
	// 	log.Errorf("handleFvmTransfer Invalid order status %d , %s", info.State, orderID)
	// 	return
	// }

	info.Value = tr.Value
	info.TxHash = tr.TxHash
	info.From = tr.From
	info.DoneHeight = getFilecoinHeight(m.cfg.LotusHTTPSAddr)

	// TODO notify
	log.Warnf("need transfer %s USDT to %s", tr.Value, info.WithdrawAddr)

	err = m.changeOrderState(types.WithdrawDone, info)
	if err != nil {
		log.Errorf("handleFvmTransfer changeOrderState %s err:%s", orderID, err.Error())
		return
	}
}

func (m *WithdrawManager) initOngeingOrders() {
	records, err := m.LoadWithdrawRecords(types.WithdrawDone)
	if err != nil {
		log.Errorln("LoadWithdrawRecords err:", err.Error())
		return
	}

	for _, info := range records {
		m.tMgr.RecoverOutstandingFvmOrders(info.To, info.OrderID)
		m.addOrder(info)
	}
}

// func (m *WithdrawManager) checkOrdersTimeout() {
// 	ticker := time.NewTicker(checkOrderInterval)
// 	defer ticker.Stop()

// 	for {
// 		<-ticker.C

// 		for _, orderRecord := range m.ongoingOrders {
// 			orderID := orderRecord.OrderID
// 			addr := orderRecord.To

// 			info, err := m.LoadWithdrawRecord(orderID)
// 			if err != nil {
// 				log.Errorf("checkOrderTimeout LoadOrderRecord %s , %s err:%s", addr, orderID, err.Error())
// 				continue
// 			}

// 			log.Debugf("checkout %s , %s ", addr, orderID)

// 			if info.State == types.ExchangeCreated && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {
// 				err := m.changeOrderState(types.ExchangeTimeout, info)
// 				if err != nil {
// 					log.Errorf("checkOrderTimeout changeOrderState %s , %s err:%s", addr, orderID, err.Error())
// 					continue
// 				}
// 			}
// 		}
// 	}
// }

func (m *WithdrawManager) getOrderIDByToAddress(to string) (string, bool) {
	for _, orderRecord := range m.ongoingOrders {
		log.Debugf("getOrderIDByToAddress orderRecord %v ", orderRecord)
		if orderRecord.To == to {
			return orderRecord.OrderID, true
		}
	}

	return "", false
}

// // CancelWithdrawOrder cancel the order
// func (m *WithdrawManager) CancelWithdrawOrder(orderID string) error {
// 	info, err := m.LoadWithdrawRecord(orderID)
// 	if err != nil {
// 		return err
// 	}

// 	if info.State != types.ExchangeCreated {
// 		return xerrors.Errorf("Invalid order status %d", info.State)
// 	}

// 	return m.changeOrderState(types.ExchangeCancel, info)
// }

func (m *WithdrawManager) changeOrderState(state types.WithdrawState, info *types.WithdrawRecord) error {
	info.DoneHeight = getFilecoinHeight(m.cfg.LotusHTTPSAddr)

	err := m.UpdateWithdrawRecord(info, state)
	if err != nil {
		return err
	}

	m.removeOrder(info.User)
	m.tMgr.RevertFvmAddress(info.To)

	return nil
}

// CreateWithdrawOrder create a withdraw order
func (m *WithdrawManager) CreateWithdrawOrder(userAddr, withdrawAddr, value string) (err error) {
	original, err := m.LoadUserToken(userAddr)
	if err != nil {
		return err
	}

	newValue, ok := utils.BigIntReduce(original, value)
	if !ok {
		return xerrors.New("Insufficient balance")
	}

	err = m.UpdateUserToken(userAddr, newValue, original)
	if err != nil {
		return err
	}

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	info := &types.WithdrawRecord{
		OrderID:       orderID,
		User:          userAddr,
		WithdrawAddr:  withdrawAddr,
		Value:         value,
		CreatedHeight: getFilecoinHeight(m.cfg.LotusHTTPSAddr),
		State:         types.WithdrawCreate,
	}

	return m.SaveWithdrawInfo(info)
}

func (m *WithdrawManager) addOrder(info *types.WithdrawRecord) error {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	if _, exist := m.ongoingOrders[info.User]; exist {
		return xerrors.New("user have order")
	}

	m.ongoingOrders[info.User] = info

	return nil
}

func (m *WithdrawManager) removeOrder(userID string) {
	m.orderLock.Lock()
	defer m.orderLock.Unlock()

	delete(m.ongoingOrders, userID)
}
