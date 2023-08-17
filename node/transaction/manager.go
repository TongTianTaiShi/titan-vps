package transaction

import (
	"sync"

	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/filecoin-project/pubsub"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("transaction")

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	notify *pubsub.PubSub
	*db.SQLDB

	cfg config.MallCfg

	tronAddrs    map[string]string
	tronAddrLock *sync.Mutex
}

// NewManager creates a new instance of the node manager
func NewManager(pb *pubsub.PubSub, getCfg dtypes.GetMallConfigFunc, db *db.SQLDB) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		notify: pb,
		cfg:    cfg,
		SQLDB:  db,

		tronAddrs:    make(map[string]string),
		tronAddrLock: &sync.Mutex{},
	}

	manager.initTronAddress(cfg.RechargeAddresses)

	go manager.watchTronTransactions()

	return manager, nil
}

// AllocateTronAddress get a fvm address
func (m *Manager) AllocateTronAddress(userID string) (string, error) {
	list, err := m.GetRechargeAddresses()
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		return "", xerrors.New("not found address")
	}

	addr := list[0]
	err = m.UpdateRechargeAddressOfUser(addr, userID)
	if err != nil {
		return "", err
	}

	m.addTronAddr(addr, userID)

	return addr, nil
}
