package config

import (
	"encoding"
	"os"
	"strconv"
	"time"
)

const (
	// RetrievalPricingDefaultMode  configures the node to use the default retrieval pricing policy.
	RetrievalPricingDefaultMode = "default"
	// RetrievalPricingExternalMode  configures the node to use the external retrieval pricing script
	// configured by the user.
	RetrievalPricingExternalMode = "external"
)

// MaxTraversalLinks configures the maximum number of links to traverse in a DAG while calculating
// CommP and traversing a DAG with graphsync; invokes a budget on DAG depth and density.
var MaxTraversalLinks uint64 = 32 * (1 << 20)

func init() {
	if envMaxTraversal, err := strconv.ParseUint(os.Getenv("TITAN_MAX_TRAVERSAL_LINKS"), 10, 64); err == nil {
		MaxTraversalLinks = envMaxTraversal
	}
}

// DefaultBasisCfg returns the default basis config
func DefaultBasisCfg() *BasisCfg {
	return &BasisCfg{
		Common: Common{
			API: API{
				ListenAddress: "0.0.0.0:5577",
			},
		},
		Timeout:               "30s",
		DryRun:                true,
		AliyunAccessKeyID:     "",
		AliyunAccessKeySecret: "",
		DatabaseAddress:       "",
	}
}

// DefaultTransactionCfg returns the default transaction config
func DefaultTransactionCfg() *TransactionCfg {
	return &TransactionCfg{
		Common: Common{
			API: API{
				ListenAddress:       "0.0.0.0:5577",
				RemoteListenAddress: "",
			},
		},
		DatabaseAddress: "",
	}
}

var (
	_ encoding.TextMarshaler   = (*Duration)(nil)
	_ encoding.TextUnmarshaler = (*Duration)(nil)
)

// Duration is a wrapper type for time.Duration
// for decoding and encoding from/to TOML
type Duration time.Duration

// UnmarshalText implements interface for TOML decoding
func (dur *Duration) UnmarshalText(text []byte) error {
	d, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*dur = Duration(d)
	return err
}

// MarshalText implements interface for TOML encoding
func (dur Duration) MarshalText() ([]byte, error) {
	d := time.Duration(dur)
	return []byte(d.String()), nil
}
