package exchange

import (
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/LMF709268224/titan-vps/node/utils"
	"github.com/filecoin-project/pubsub"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("exchange")

// RechargeManager manager recharge order
type RechargeManager struct {
	*db.SQLDB
	cfg          config.MallCfg
	notification *pubsub.PubSub

	tMgr *transaction.Manager
}

// NewRechargeManager returns a new manager instance
func NewRechargeManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, fm *transaction.Manager) (*RechargeManager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &RechargeManager{
		SQLDB:        sdb,
		notification: pb,
		cfg:          cfg,

		tMgr: fm,
	}

	go m.subscribeEvents()

	return m, nil
}

func (m *RechargeManager) subscribeEvents() {
	subTransfer := m.notification.Sub(types.EventTronTransferWatch.String())
	defer m.notification.Unsub(subTransfer)

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

	exist, err := m.RechargeRecordExists(tr.TxHash)
	if err != nil {
		log.Errorf("RechargeRecordExists:%v", err)
		return
	}

	if exist {
		return
	}

	userID := tr.UserID
	original, err := m.LoadUserBalance(userID)
	if err != nil {
		log.Errorf("%s LoadUserToken err:%s", userID, err.Error())
		return
	}

	info := &types.RechargeRecord{
		OrderID: tr.TxHash,
		UserID:  userID,
		Value:   tr.Value,
		From:    tr.From,
		State:   types.RechargeDone,
		To:      tr.To,
	}

	value := utils.BigIntAdd(original, tr.Value)

	err = m.SaveRechargeRecordAndUserBalance(info, value, original)
	if err != nil {
		log.Errorf("%s SaveRechargeRecordAndUserBalance err:%s", info.OrderID, err.Error())
	}
}
