package exchange

import (
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/trxbridge"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/api"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/hexutil"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/pubsub"
	"github.com/google/uuid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/smirkcat/hdwallet"
	"golang.org/x/xerrors"
	"google.golang.org/protobuf/proto"
)

var log = logging.Logger("exchange")

const (
	checkOrderInterval = 10 * time.Second
	checkBlockInterval = 10 * time.Second
	checkLimit         = 100
	orderTimeoutTime   = 5 * time.Minute
)

// Manager manager order
type Manager struct {
	*db.SQLDB
	cfg    config.BasisCfg
	notify *pubsub.PubSub

	rMgr   *RechargeMgr
	height int64

	usabilityAddrs map[string]string
	usedAddrs      map[string]string
	addrLock       *sync.Mutex

	ongoingOrders map[string]string
	orderLock     *sync.Mutex
}

// NewManager returns a new manager instance
func NewManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetBasisConfigFunc) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	rMgr, err := NewRechargeMgr(sdb, pb, cfg)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		SQLDB:  sdb,
		notify: pb,
		cfg:    cfg,

		rMgr:   rMgr,
		height: 39130320,

		usabilityAddrs: make(map[string]string),
		usedAddrs:      make(map[string]string),
		addrLock:       &sync.Mutex{},

		ongoingOrders: make(map[string]string),
		orderLock:     &sync.Mutex{},
	}

	m.initPaymentAddress(m.cfg.RechargeAddress)

	go m.watchTransactions()
	go m.checkOrderTimeout()

	return m, nil
}

func (m *Manager) initPaymentAddress(as []string) {
	m.addrLock.Lock()
	defer m.addrLock.Unlock()

	for _, addr := range as {
		m.usabilityAddrs[addr] = ""
	}
}

// GetGrpcClient
func (m *Manager) getGrpcClient() (*trxbridge.GrpcClient, error) {
	node := trxbridge.NewGrpcClient(m.cfg.TrxHTTPSAddr)
	err := node.Start()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (m *Manager) watchTransactions() {
	ticker := time.NewTicker(checkBlockInterval)
	defer ticker.Stop()

	client, err := m.getGrpcClient()
	if err != nil {
		log.Errorln("getGrpcClient err :", err.Error())
		return
	}

	for {
		<-ticker.C

		// blockExtention, err := client.GetBlockByNum(m.height)
		// if err == nil {
		// 	err = m.handleBlock(blockExtention)
		// 	if err == nil {
		// 		m.height++
		// 	}
		// }

		if len(m.usedAddrs) > 0 {
			blocks, err := client.GetBlockByLatestNum(4)
			if err == nil && len(blocks.Block) > 0 {
				m.handleBlocks(blocks)
			}
		}
	}
}

func (m *Manager) handleBlocks(blocks *api.BlockListExtention) {
	for _, v := range blocks.Block {
		err := m.handleBlock(v)
		if err != nil {
			log.Errorln(" handleBlock err :", err.Error())
		}
	}
}

func (m *Manager) handleBlock(blockExtention *api.BlockExtention) error {
	if blockExtention == nil || blockExtention.BlockHeader == nil {
		return xerrors.New("block is nil")
	}

	bNum := blockExtention.BlockHeader.RawData.Number
	log.Infoln(" handleBlock bNum :", bNum)

	bid := hexutil.Encode(blockExtention.Blockid)

	for _, te := range blockExtention.Transactions {
		if len(te.Transaction.GetRet()) == 0 {
			continue
		}

		state := te.Transaction.GetRet()[0].ContractRet
		txid := hexutil.Encode(te.Txid)

		for _, contract := range te.Transaction.RawData.Contract {
			m.filterTransaction(contract, txid, bid, bNum, state)
		}
	}

	return nil
}

func (m *Manager) filterTransaction(contract *core.Transaction_Contract, txid, bid string, bNum int64, state core.Transaction_ResultContractResult) {
	if contract.Type == core.Transaction_Contract_TriggerSmartContract {
		// trc20
		unObj := &core.TriggerSmartContract{}
		err := proto.Unmarshal(contract.Parameter.GetValue(), unObj)
		if err != nil {
			log.Errorf("parse trc20 err: %v", err)
			return
		}

		contractAddress := hdwallet.EncodeCheck(unObj.GetContractAddress())

		if contractAddress != m.cfg.TrxContractorAddr {
			return
		}

		from := hdwallet.EncodeCheck(unObj.GetOwnerAddress())
		data := unObj.GetData()

		to, amount, isOk := m.decodeData(data)
		if !isOk {
			return
		}

		m.handleTransfer(txid, from, to, bid, bNum, amount, contractAddress, state)
	}
}

func (m *Manager) handleTransfer(mCid, from, to, blockCid string, height int64, amount string, contract string, state core.Transaction_ResultContractResult) {
	log.Infof("检查地址:%s,%s,%s,%s,%s,%s", mCid, to, from, contract, amount, state)

	if hash, exist := m.usedAddrs[to]; exist {
		info, err := m.LoadRechargeRecord(hash)
		if err != nil {
			log.Errorf("handleTransfer LoadOrderRecord %s , %s err:%s", to, hash, err.Error())
			return
		}
		info.DoneState = state
		info.Value = amount
		info.TxHash = mCid
		info.From = from

		err = m.UpdateRechargeRecord(info, types.RechargeDone)
		if err != nil {
			log.Errorf("handleTransfer UpdateRechargeRecord %s , %s err:%s", to, hash, err.Error())
			return
		}

		m.removeOrder(info.User)
		m.revertPayeeAddress(info.To)
	}
}

func (m *Manager) checkOrderTimeout() {
	ticker := time.NewTicker(checkOrderInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		for addr, hash := range m.usedAddrs {
			info, err := m.LoadRechargeRecord(hash)
			if err != nil {
				log.Errorf("checkOrderTimeout LoadOrderRecord %s , %s err:%s", addr, hash, err.Error())
				continue
			}

			log.Debugf("checkout %s , %s ", addr, hash)

			if info.State == types.RechargeCreated && info.CreatedTime.Add(orderTimeoutTime).Before(time.Now()) {
				err := m.UpdateRechargeRecord(info, types.RechargeTimeout)
				if err != nil {
					log.Errorf("checkOrderTimeout UpdateRechargeRecord %s , %s err:%s", addr, hash, err.Error())
					continue
				}

				m.removeOrder(info.User)
				m.revertPayeeAddress(info.To)
			}
		}
	}
}

func (m *Manager) addressExist(addr string) bool {
	if _, exist := m.usabilityAddrs[addr]; exist {
		return true
	}

	if _, exist := m.usedAddrs[addr]; exist {
		return true
	}

	return false
}

func (m *Manager) decodeData(trc20 []byte) (to string, amount string, flag bool) {
	if len(trc20) >= 68 {
		if hexutil.Encode(trc20[:4]) != "a9059cbb" {
			return
		}
		trc20[15] = 65 // 0x41

		bb := common.TrimLeftZeroes(trc20[36:68])
		bu := new(big.Int).SetBytes(bb)
		amount = bu.String()

		to = hdwallet.EncodeCheck(trc20[15:36])
		flag = true
	}
	return
}

func (m *Manager) CreateRechargeOrder(userAddr, rechargeAddr string) (string, error) {
	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	err := m.addOrder(userAddr, orderID)
	if err != nil {
		return "", err
	}

	addr, err := m.allocatePayeeAddress(orderID)
	if err != nil {
		return "", err
	}

	info := &types.RechargeRecord{
		ID:           orderID,
		User:         userAddr,
		To:           addr,
		RechargeAddr: rechargeAddr,
		// State: ,
	}

	err = m.SaveRechargeInfo(info)
	if err != nil {
		return "", err
	}

	return addr, nil
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
