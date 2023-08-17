package dtypes

import (
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/ipfs/go-datastore"
	"github.com/multiformats/go-multiaddr"
)

// MetadataDS stores metadata.
type MetadataDS datastore.Batching

type APIAlg jwt.HMACSHA

type APIEndpoint multiaddr.Multiaddr

// InternalIP local network address
type InternalIP string

type (
	NodeMetadataPath string
	AssetsPaths      []string
)

// SetTransactionConfigFunc is a function which is used to
// sets the transaction config.
type SetTransactionConfigFunc func(cfg config.TransactionCfg) error

// GetTransactionConfigFunc is a function which is used to
// get the sealing config.
type GetTransactionConfigFunc func() (config.TransactionCfg, error)

// SetMallConfigFunc is a function which is used to
// sets the mall config.
type SetMallConfigFunc func(cfg config.MallCfg) error

// GetMallConfigFunc is a function which is used to
// get the sealing config.
type GetMallConfigFunc func() (config.MallCfg, error)
