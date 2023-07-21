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

// SetBasisConfigFunc is a function which is used to
// sets the basis config.
type SetBasisConfigFunc func(cfg config.BasisCfg) error

// GetBasisConfigFunc is a function which is used to
// get the sealing config.
type GetBasisConfigFunc func() (config.BasisCfg, error)
