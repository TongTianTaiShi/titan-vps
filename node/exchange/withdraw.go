package exchange

import (
	"strings"

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
	cfg          config.MallCfg
	notification *pubsub.PubSub

	tMgr *transaction.Manager
}

// NewWithdrawManager returns a new manager instance
func NewWithdrawManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, fm *transaction.Manager) (*WithdrawManager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &WithdrawManager{
		SQLDB:        sdb,
		notification: pb,
		cfg:          cfg,

		tMgr: fm,
	}

	return m, nil
}

// CreateWithdrawOrder create a withdraw order
func (m *WithdrawManager) CreateWithdrawOrder(userID, withdrawAddr, value string) (err error) {
	original, err := m.LoadUserBalance(userID)
	if err != nil {
		return err
	}

	newValue, ok := utils.BigIntReduce(original, value)
	if !ok {
		return xerrors.New("Insufficient balance")
	}

	err = m.UpdateUserBalance(userID, newValue, original)
	if err != nil {
		return err
	}

	hash := uuid.NewString()
	orderID := strings.Replace(hash, "-", "", -1)

	info := &types.WithdrawRecord{
		OrderID:       orderID,
		User:          userID,
		WithdrawAddr:  withdrawAddr,
		Value:         value,
		CreatedHeight: getFilecoinHeight(m.cfg.LotusHTTPSAddr),
		State:         types.WithdrawCreate,
	}

	return m.SaveWithdrawInfo(info)
}
