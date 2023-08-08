package types

import (
	"time"

	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"
)

type Hellos struct {
	Msg string
}

// User user info
type User struct {
	UUID      string    `db:"uuid" json:"uuid"`
	UserName  string    `db:"user_name" json:"user_name"`
	PassHash  string    `db:"pass_hash" json:"pass_hash"`
	Address   string    `db:"address" json:"address"`
	Public    string    `db:"public" json:"public"`
	Token     string    `db:"token" json:"token"`
	Role      int32     `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DescribePriceResponse struct {
	Currency      string
	OriginalPrice float32
	TradePrice    float32
}

// todo
type CreateInstanceReq struct {
	Id                      string `db:"id"`
	RegionId                string `db:"region_id"`
	InstanceType            string `db:"instance_type"`
	DryRun                  bool   `db:"dry_run"`
	ImageId                 string `db:"image_id"`
	SecurityGroupId         string `db:"security_group_id"`
	InstanceChargeType      string `db:"instanceCharge_type"`
	PeriodUnit              string `db:"period_unit"`
	Period                  int32  `db:"period"`
	InternetMaxBandwidthOut int32  `db:"bandwidth_out"`
	InternetMaxBandwidthIn  int32  `db:"bandwidth_in"`
}

type CreateInstanceResponse struct {
	InstanceID       string  `db:"instance_id"`
	OrderId          string  `db:"order_id"`
	RequestId        string  `db:"request_id"`
	TradePrice       float32 `db:"trade_price"`
	PublicIpAddress  string  `db:"public_ip_address"`
	PrivateKey       string
	PrivateKeyStatus int `db:"private_key_status"`
}

type CreateKeyPairResponse struct {
	KeyPairID      string
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
	User          string    `db:"user_addr"`
	To            string    `db:"to_addr"`
	Value         int64     `db:"value"`
	State         int64     `db:"state"`
	DoneState     int64     `db:"done_state"`
	CreatedHeight int64     `db:"created_height"`
	CreatedTime   time.Time `db:"created_time"`
	DoneTime      time.Time `db:"done_time"`
	DoneHeight    int64     `db:"done_height"`
	VpsID         int64     `db:"vps_id"`
	Msg           string    `db:"msg"`
	TxHash        string    `db:"tx_hash"`
}

// RechargeState
type RechargeState int64

// Constants defining various states of the recharge process.
const (
	// RechargeCreated recharge order
	RechargeCreated RechargeState = iota
	// RechargeDone the order done
	RechargeDone
	// RechargeTimeout order
	RechargeTimeout
)

// RechargeRecord represents information about an recharge record
type RechargeRecord struct {
	ID            string                                `db:"id"`
	From          string                                `db:"from_addr"`
	User          string                                `db:"user_addr"`
	To            string                                `db:"to_addr"`
	Value         string                                `db:"value"`
	State         RechargeState                         `db:"state"`
	DoneState     core.Transaction_ResultContractResult `db:"done_state"`
	CreatedHeight int64                                 `db:"created_height"`
	CreatedTime   time.Time                             `db:"created_time"`
	DoneTime      time.Time                             `db:"done_time"`
	Msg           string                                `db:"msg"`
	DoneHeight    int64                                 `db:"done_height"`
	TxHash        string                                `db:"tx_hash"`
	RechargeAddr  string                                `db:"recharge_addr"`
	RechargeHash  string                                `db:"recharge_hash"`
}

type CreateOrderReq struct {
	Vps  CreateInstanceReq
	User string
}

type PaymentCompletedReq struct {
	OrderID       string
	TransactionID string
}
type UserReq struct {
	UserId    string
	Signature string
	Address   string
	PublicKey string
	Token     string
}
type UserResponse struct {
	UserId   string
	SignCode string
	Token    string
}

type UserInfoTmp struct {
	UserLogin UserResponse
	OrderInfo OrderRecord
}
type Token struct {
	TokenString string
	UserId      string
	Expiration  time.Time
}

// EventTopics represents topics for pub/sub events
type EventTopics string

const (
	// EventTransferWatch node online event
	EventTransferWatch EventTopics = "transfer_watch"
	// EventTransferReq request transfer
	EventTransferReq EventTopics = "transfer_req"
	// EventTransferRep reply transfer
	EventTransferRep EventTopics = "transfer_rep"
)

func (t EventTopics) String() string {
	return string(t)
}

type FvmTransferWatch struct {
	ID    string
	From  string
	To    string
	Value int64
}

type FvmTransferReq struct {
	ID    string
	To    string
	Value string
}

type FvmTransferRep struct {
	ID     string
	TxHash string
	Msg    string
}
