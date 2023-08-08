package exchange

import (
	"sync"

	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/filecoin-project/pubsub"
)

type RechargeMgr struct {
	*db.SQLDB
	cfg    config.BasisCfg
	notify *pubsub.PubSub

	ongoingRecharges map[string]string
	rechargeLock     *sync.Mutex

	usabilityAddrs map[string]string
	usedAddrs      map[string]string
	addrLock       *sync.Mutex
}

// NewRechargeMgr returns a new manager instance
func NewRechargeMgr(sdb *db.SQLDB, pb *pubsub.PubSub, cfg config.BasisCfg) (*RechargeMgr, error) {
	m := &RechargeMgr{
		SQLDB:            sdb,
		notify:           pb,
		ongoingRecharges: make(map[string]string),
		rechargeLock:     &sync.Mutex{},
		cfg:              cfg,

		usabilityAddrs: make(map[string]string),
		usedAddrs:      make(map[string]string),
		addrLock:       &sync.Mutex{},
	}

	return m, nil
}

