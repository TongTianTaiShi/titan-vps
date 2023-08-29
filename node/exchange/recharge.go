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

// RechargeManager manages recharge orders
type RechargeManager struct {
	*db.SQLDB
	cfg          config.MallCfg
	notification *pubsub.PubSub
	tMgr         *transaction.Manager
}

// NewRechargeManager creates a new manager instance for handling recharge orders
func NewRechargeManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, fm *transaction.Manager) (*RechargeManager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &RechargeManager{
		SQLDB:        sdb,
		notification: pb,
		cfg:          cfg,
		tMgr:         fm,
	}

	go m.subscribeEvents()

	return m, nil
}

// subscribeEvents subscribes to Tron transfer events and processes them
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

// handleTronTransfer handles Tron transfer events and updates user balances
func (m *RechargeManager) handleTronTransfer(tr *types.TronTransferWatch) {
	if tr.State != core.Transaction_Result_SUCCESS {
		// If the transaction state is not successful, skip processing.
		return
	}

	// Check if the recharge record already exists for this transaction.
	exist, err := m.RechargeRecordExists(tr.TxHash)
	if err != nil {
		log.Errorf("RechargeRecordExists error: %v", err)
		return
	}

	if exist {
		// If the recharge record already exists, skip processing.
		return
	}

	userID := tr.UserID

	// Load the original user balance.
	original, err := m.LoadUserBalance(userID)
	if err != nil {
		log.Errorf("%s LoadUserBalance error: %s", userID, err.Error())
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

	// Calculate the new user balance.
	value, err := utils.AddBigInt(original, tr.Value)
	if err != nil {
		log.Errorf("%s BigIntAdd error: %s", userID, err.Error())
		return
	}

	// Save the recharge record and update the user balance.
	err = m.SaveRechargeRecordAndUserBalance(info, value, original)
	if err != nil {
		log.Errorf("%s SaveRechargeRecordAndUserBalance error: %s", info.OrderID, err.Error())
	}
}
