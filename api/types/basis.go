package types

import (
	"time"
)

type Hellos struct {
	Msg string
}

// User user info
type User struct {
	UUID     string `db:"uuid" json:"uuid"`
	UserName string `db:"user_name" json:"user_name"`
	PassHash string `db:"pass_hash" json:"pass_hash"`
	// UserEmail string    `db:"user_email" json:"user_email"`
	// Role      int32     `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DescribePriceResponse struct {
	Currency      string
	OriginalPrice float32
	TradePrice    float32
}

type CreateInstanceResponse struct {
	InstanceId      string
	OrderId         string
	RequestId       string
	TradePrice      float32
	PublicIpAddress string
	PrivateKey      string
}
type CreateKeyPairResponse struct {
	KeyPairId      string
	KeyPairName    string
	PrivateKeyBody string
}
type AttachKeyPairResponse struct {
	Code       string
	InstanceId string
	Message    string
	Success    string
}

// OrderRecord represents information about an asset record
type OrderRecord struct {
	CID                   string    `db:"cid"`
	Hash                  string    `db:"hash"`
	NeedEdgeReplica       int64     `db:"edge_replicas"`
	TotalSize             int64     `db:"total_size"`
	TotalBlocks           int64     `db:"total_blocks"`
	Expiration            time.Time `db:"expiration"`
	CreatedTime           time.Time `db:"created_time"`
	EndTime               time.Time `db:"end_time"`
	NeedCandidateReplicas int64     `db:"candidate_replicas"`
	State                 string    `db:"state"`
	NeedBandwidth         int64     `db:"bandwidth"` // unit:MiB/s

	RetryCount        int64 `db:"retry_count"`
	ReplenishReplicas int64 `db:"replenish_replicas"`

	SPCount int64
}
