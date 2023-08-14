package transaction

import (
	"sync"

	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/LMF709268224/titan-vps/node/modules/dtypes"
	"github.com/filecoin-project/pubsub"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("transaction")

// Manager is the node manager responsible for managing the online nodes
type Manager struct {
	notify *pubsub.PubSub
	*db.SQLDB

	cfg config.BasisCfg

	usabilityFvmAddrs map[string]string
	usedFvmAddrs      map[string]string
	fvmAddrLock       *sync.Mutex

	tronAddr string

	addrWait sync.WaitGroup
}

// NewManager creates a new instance of the node manager
func NewManager(pb *pubsub.PubSub, getCfg dtypes.GetBasisConfigFunc, db *db.SQLDB) (*Manager, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		notify: pb,
		cfg:    cfg,
		SQLDB:  db,

		usabilityFvmAddrs: make(map[string]string),
		usedFvmAddrs:      make(map[string]string),
		fvmAddrLock:       &sync.Mutex{},
	}

	manager.addrWait.Add(2)
	manager.initFvmAddress(cfg.PaymentAddresses)

	manager.tronAddr = cfg.RechargeAddress

	go manager.watchFvmTransactions()
	go manager.watchTronTransactions()

	return manager, nil
}
