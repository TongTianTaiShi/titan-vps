package exchange

import (
	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/api/terrors"
	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/LMF709268224/titan-vps/node/transaction"
	"github.com/LMF709268224/titan-vps/node/utils"
	"github.com/filecoin-project/pubsub"
	"github.com/google/uuid"
)

// WithdrawManager manages withdrawal orders
type WithdrawManager struct {
	*db.SQLDB
	cfg          config.MallCfg
	notification *pubsub.PubSub
	tMgr         *transaction.Manager
}

// NewWithdrawManager creates a new manager instance for handling withdrawal orders
func NewWithdrawManager(sdb *db.SQLDB, pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, fm *transaction.Manager) (*WithdrawManager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	m := &WithdrawManager{
		SQLDB:        sdb,
		notification: pb,
		cfg:          cfg,
		tMgr:         fm,
	}

	return m, nil
}

// CreateWithdrawOrder creates a withdrawal order for a user
func (m *WithdrawManager) CreateWithdrawOrder(userID, withdrawAddr, value string) error {
	// Load the original user balance.
	original, err := m.LoadUserBalance(userID)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	// Reduce the user's balance to cover the withdrawal amount.
	newValue, err := utils.ReduceBigInt(original, value)
	if err != nil {
		return &api.ErrWeb{Code: terrors.InsufficientBalance.Int(), Message: terrors.InsufficientBalance.String()}
	}

	// Generate a unique order ID.
	orderID := uuid.NewString()

	// Create a WithdrawRecord to store withdrawal information.
	info := &types.WithdrawRecord{
		OrderID:      orderID,
		UserID:       userID,
		WithdrawAddr: withdrawAddr,
		Value:        value,
		State:        types.WithdrawCreate,
	}

	// Save the withdrawal record and update the user balance.
	err = m.SaveWithdrawInfoAndUserBalance(info, newValue, original)
	if err != nil {
		return &api.ErrWeb{Code: terrors.DatabaseError.Int(), Message: err.Error()}
	}

	return nil
}
