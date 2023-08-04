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
	Role      int32     `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DescribePriceResponse struct {
	Currency      string
	OriginalPrice float32
	TradePrice    float32
}

type CreateInstanceReq struct {
	RegionId                string
	InstanceType            string
	DryRun                  bool
	ImageId                 string
	SecurityGroupId         string
	InstanceChargeType      string
	PeriodUnit              string
	Period                  int32
	InternetMaxBandwidthOut int32
	InternetMaxBandwidthIn  int32
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

// OrderRecord represents information about an order record
type OrderRecord struct {
	OrderID       string    `db:"order_id"`
	From          string    `db:"from_addr"`
	To            string    `db:"to_addr"`
	Value         int64     `db:"value"`
	State         int64     `db:"state"`
	DoneState     int64     `db:"done_state"`
	CreatedHeight int64     `db:"created_height"`
	CreatedTime   time.Time `db:"created_time"`
	DoneTime      time.Time `db:"done_time"`
	DoneHeight    int64     `db:"done_height"`
	VpsID         string    `db:"vps_id"`
}

type CreateOrderReq struct {
	Vps  CreateInstanceReq
	User string
}

type PaymentCompletedReq struct {
	OrderID       string
	TransactionID string
}

// EventTopics represents topics for pub/sub events
type EventTopics string

const (
	// EventTransfer node online event
	EventTransfer EventTopics = "transfer"
)

func (t EventTopics) String() string {
	return string(t)
}

type FvmTransfer struct {
	From  string
	To    string
	Value int64
}
