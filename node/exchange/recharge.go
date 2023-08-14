package exchange

import (
	"strings"
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

	tMgr *transaction.Manager
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
	}

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
	if tr.State != core.Transaction_Result_SUCCESS {
		return
	}

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	userAddr := tr.UserAddr
	height := getTronHeight(m.cfg.TrxHTTPSAddr)

	info := &types.RechargeRecord{
		OrderID:       orderID,
		User:          userAddr,
		CreatedHeight: height,
		DoneHeight:    height,
		Value:         tr.Value,
		TxHash:        tr.TxHash,
		From:          tr.From,
	}

	info.Msg = tr.State.String()

	client := filecoinbridge.NewGrpcClient(m.cfg.LotusHTTPSAddr, m.cfg.TitanContractorAddr)
	hash, err := client.Mint(m.cfg.PrivateKeyStr, userAddr, tr.Value)
	if err != nil {
		info.Msg = err.Error()
		info.State = types.RechargeRefund
	} else {
		info.RechargeHash = hash
		info.State = types.RechargeDone
	}

	m.SaveRechargeInfo(info)
}
