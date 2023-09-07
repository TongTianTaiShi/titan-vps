package types

import (
	"time"

	"github.com/LMF709268224/titan-vps/lib/trxbridge/core"
)

// OrderState represents the state of an order in the process .
type OrderState int64

// Constants defining various states of the order process.
const (
	// Created order
	Created OrderState = iota
	// WaitingPayment Waiting for user to payment order
	WaitingPayment
	// BuyGoods buy goods
	BuyGoods
	// Done the order done
	Done
)

// Int returns the int representation of the order state.
func (s OrderState) Int() int64 {
	return int64(s)
}

// OrderDoneState represents the different states of a completed order.
type OrderDoneState int64

// Constants defining various states of a completed order.
const (
	// OrderDoneStateSuccess represents the state when the order is completed successfully.
	OrderDoneStateSuccess OrderDoneState = iota
	// OrderDoneStateTimeout represents the state when the order has timed out.
	OrderDoneStateTimeout
	// OrderDoneStateCancel represents the state when the user has canceled the order.
	OrderDoneStateCancel
	// OrderDoneStatePurchaseFailed represents the state when the purchase has failed.
	OrderDoneStatePurchaseFailed
)

// String returns the string representation of the order done state.
func (s OrderDoneState) String() string {
	switch s {
	case 0:
		return "Success"
	case 1:
		return "Timeout"
	case 2:
		return "Cancel"
	case 3:
		return "PurchaseFailed"
	}

	return "Not found"
}

// Int returns the int representation of the order done state.
func (s OrderDoneState) Int() int64 {
	return int64(s)
}

// OrderType order type
type OrderType int64

// Constants defining various states of the order process.
const (
	// BuyVPS order
	BuyVPS OrderType = iota
	// Renew order
	RenewVPS
)

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

type DescribePriceReq struct {
	RegionId                     string
	InstanceType                 string
	PriceUnit                    string
	ImageID                      string
	InternetChargeType           string
	SystemDiskCategory           string
	SystemDiskSize               int32
	Period                       int32
	Amount                       int32
	InternetMaxBandwidthOut      int32
	DescribePriceRequestDataDisk []DescribePriceRequestDataDisk
}

type DescribePriceRequestDataDisk struct {
	Category         string
	PerformanceLevel string
	Size             int64
}

type DescribePriceResponse struct {
	Currency      string
	OriginalPrice float32
	TradePrice    float32
	USDPrice      float32
}

type DescribeImageResponse struct {
	ImageId      string
	ImageName    string
	ImageFamily  string
	Platform     string
	OSType       string
	OSName       string
	Architecture string
}

type CreateOrderReq struct {
	CreateInstanceReq
	Amount int32
}

type RenewOrderReq struct {
	InstanceId string `db:"instance_id"`
	PeriodUnit string `db:"period_unit"`
	Period     int32  `db:"period"`
	Renew      int    `db:"renew"`
}

type SetRenewOrderReq struct {
	RegionID   string `db:"region_id"`
	InstanceId string `db:"instance_id"`
	PeriodUnit string `db:"period_unit"`
	Period     int32  `db:"period"`
	Renew      int    `db:"renew"`
}

type CreateInstanceResponse struct {
	InstanceID       string  `db:"instance_id"`
	OrderId          string  `db:"order_id"`
	RequestId        string  `db:"request_id"`
	TradePrice       float32 `db:"trade_price"`
	PublicIpAddress  string  `db:"public_ip_address"`
	PrivateKeyStatus int     `db:"private_key_status"`
	PrivateKey       string
	AccessKey        string
}

type DescribeInstanceTypeReq struct {
	RegionId         string
	MemorySize       float32
	CpuArchitecture  string
	InstanceCategory string
	CpuCoreCount     int32
	MaxResults       int64
	NextToken        string
}

type AvailableResourceReq struct {
	RegionId            string
	DestinationResource string
	InstanceChargeType  string
	InstanceType        string
	ResourceType        string
}

type RenewInstanceRequest struct {
	RegionId   string
	InstanceId string
	PeriodUnit string
	Period     int32
}

type AvailableResourceResponse struct {
	Min   int32
	Max   int32
	Value string
	Unit  string
}

type DescribeRecommendInstanceTypeReq struct {
	RegionId           string
	Memory             float32
	Cores              int32
	InstanceChargeType string
}

type DescribeRecommendInstanceResponse struct {
	Memory             int32
	Cores              int32
	InstanceType       string
	InstanceTypeFamily string
}

type DescribeInstanceTypeResponse struct {
	InstanceTypes []*DescribeInstanceType
	NextToken     string
}

type DescribeInstanceType struct {
	InstanceTypeId         string
	MemorySize             float32
	CpuArchitecture        string
	InstanceCategory       string
	CpuCoreCount           int32
	AvailableZone          int
	InstanceTypeFamily     string
	PhysicalProcessorModel string
	NextToken              string
	Status                 string
}

type InstanceTypeFromBaseReq struct {
	RegionId            string
	MemorySize          float32
	CpuArchitecture     string
	InstanceCategory    string
	CpuCoreCount        int32
	Limit, Page, Offset int64
}

type InstanceTypeResponse struct {
	List  []*DescribeInstanceTypeFromBase
	Total int
}

type DescribeInstanceTypeFromBase struct {
	RegionId               string    `db:"region_id"`
	InstanceTypeId         string    `db:"instance_type_id"`
	MemorySize             float32   `db:"memory_size"`
	CpuArchitecture        string    `db:"cpu_architecture"`
	InstanceCategory       string    `db:"instance_category"`
	CpuCoreCount           int32     `db:"cpu_core_count"`
	AvailableZone          int       `db:"available_zone"`
	InstanceTypeFamily     string    `db:"instance_type_family"`
	PhysicalProcessorModel string    `db:"physical_processor_model"`
	OriginalPrice          float32   `db:"original_price"`
	Price                  float32   `db:"price"`
	Status                 string    `db:"status"`
	CreatedTime            time.Time `db:"created_time"`
	UpdatedTime            time.Time `db:"updated_time"`
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
	OrderID     string         `db:"order_id"`
	UserID      string         `db:"user_id"`
	Value       string         `db:"value"`
	State       OrderState     `db:"state"`
	DoneState   OrderDoneState `db:"done_state"`
	CreatedTime time.Time      `db:"created_time"`
	DoneTime    time.Time      `db:"done_time"`
	VpsID       int64          `db:"vps_id"`
	Msg         string         `db:"msg"`
	CycleTime   string         `db:"cycle_time"`
	Expiration  time.Time      `db:"expiration"`
	OrderType   OrderType      `db:"order_type"`
}

type OrderRecordResponse struct {
	Total int
	List  []*OrderRecord
}

// RechargeState Recharge order state
type RechargeState int64

// Constants defining various states of the recharge process.
const (
	// RechargeCreate Recharge create
	RechargeCreate RechargeState = iota
	// RechargeDone Recharge done
	RechargeDone
	// RechargeRefund Recharge Refund
	RechargeRefund
)

// WithdrawState Withdraw order state
type WithdrawState int64

// Constants defining various states of the recharge process.
const (
	// WithdrawCreate Withdraw create
	WithdrawCreate WithdrawState = iota
	// WithdrawDone Withdraw done
	WithdrawDone
	// WithdrawRefund Withdraw Refund
	WithdrawRefund
)

type LoginType int64

// Constants defining various states of the recharge process.
const (
	// LoginTypeMetaMask
	LoginTypeMetaMask LoginType = iota
	// LoginTypeTron
	LoginTypeTron

	LoginTypeEmail
	LoginTypeFilecoin
)

func (l LoginType) String() string {
	switch l {
	case LoginTypeMetaMask:
		return "MetaMask"
	case LoginTypeTron:
		return "Tron"
	case LoginTypeEmail:
		return "Email"
	case LoginTypeFilecoin:
		return "Filecoin"
	default:
		return "Not Found"
	}
}

type RechargeResponse struct {
	Total int
	List  []*RechargeRecord
}

// RechargeRecord represents information about an recharge record
type RechargeRecord struct {
	OrderID     string        `db:"order_id"`
	From        string        `db:"from_addr"`
	UserID      string        `db:"user_id"`
	To          string        `db:"to_addr"`
	Value       string        `db:"value"`
	State       RechargeState `db:"state"`
	CreatedTime time.Time     `db:"created_time"`
	DoneTime    time.Time     `db:"done_time"`
}

type GetWithdrawRequest struct {
	Limit     int64
	Offset    int64
	StartDate string
	EndDate   string
	State     string
	UserID    string
}

type GetWithdrawResponse struct {
	Total int
	List  []*WithdrawRecord
}

// WithdrawRecord represents information about an withdraw record
type WithdrawRecord struct {
	OrderID      string        `db:"order_id"`
	UserID       string        `db:"user_id"`
	Value        string        `db:"value"`
	State        WithdrawState `db:"state"`
	CreatedTime  time.Time     `db:"created_time"`
	DoneTime     time.Time     `db:"done_time"`
	WithdrawAddr string        `db:"withdraw_addr"`
	WithdrawHash string        `db:"withdraw_hash"`
	Executor     string        `db:"executor"`
}

type PaymentCompletedReq struct {
	OrderID       string
	TransactionID string
}

type UserReq struct {
	UserId    string
	Signature string
	Type      LoginType
}

type LoginResponse struct {
	UserId   string
	SignCode string
	Token    string
}

type UserInfoTmp struct {
	UserLogin LoginResponse
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
	// EventFvmTransferWatch node online event
	EventFvmTransferWatch EventTopics = "fvm_transfer_watch"
	// EventTronTransferWatch node online event
	EventTronTransferWatch EventTopics = "tron_transfer_watch"
)

func (t EventTopics) String() string {
	return string(t)
}

type FvmTransferWatch struct {
	TxHash string
	From   string
	To     string
	Value  string
}

type TronTransferWatch struct {
	TxHash string
	From   string
	To     string
	Value  string
	State  core.Transaction_ResultContractResult
	UserID string
}

type RechargeAddress struct {
	Addr   string `db:"addr"`
	UserID string `db:"user_id"`
}

type GetInstanceResponse struct {
	Total int
	List  []*InstanceDetails
}

type InstanceDefault struct {
	InstanceType string  `db:"instance_type"`
	RegionId     string  `db:"region_id"`
	Price        float32 `db:"price"`
}

type InstanceDetails struct {
	ID                 int64     `db:"id"`
	InstanceId         string    `db:"instance_id"`
	InstanceName       string    `db:"instance_name"`
	RegionId           string    `db:"region_id"`
	UserID             string    `db:"user_id"`
	Memory             float32   `db:"memory"`
	MemoryUsed         float32   `db:"memory_used"`
	Cores              int32     `db:"cores"`
	CoresUsed          float32   `db:"cores_used"`
	OSType             string    `db:"os_type"`
	OrderID            string    `db:"order_id"`
	InstanceType       string    `db:"instance_type"`
	ImageID            string    `db:"image_id"`
	SecurityGroupId    string    `db:"security_group_id"`
	InstanceChargeType string    `db:"instance_charge_type"`
	InternetChargeType string    `db:"internet_charge_type"`
	BandwidthOut       int32     `db:"bandwidth_out"`
	BandwidthIn        int32     `db:"bandwidth_in"`
	SystemDiskSize     int32     `db:"system_disk_size"`
	IpAddress          string    `db:"ip_address"`
	SystemDiskCategory string    `db:"system_disk_category"`
	DataDiskString     string    `db:"data_disk"`
	AutoRenew          int       `db:"auto_renew"`
	PeriodUnit         string    `db:"period_unit"`
	Period             int32     `db:"period"`
	Value              string    `db:"value"`
	CreatedTime        time.Time `db:"created_time"`
	ExpiredTime        string    `db:"expired_time"`
	AccessKey          string    `db:"access_key"`
	State              string    `db:"state"`
	Renew              string    `db:"renew"`
	DataDisk           []DescribePriceRequestDataDisk
	Executor           string    `db:"executor"`
	RefundTime         string    `db:"refund_time"`
	UpdateTime         time.Time `db:"update_time"`
}

type CreateInstanceReq struct {
	RegionId                string `db:"region_id"`
	InstanceType            string `db:"instance_type"`
	ImageID                 string `db:"image_id"`
	InternetChargeType      string `db:"internet_charge_type"`
	PeriodUnit              string `db:"period_unit"`
	Period                  int32  `db:"period"`
	InternetMaxBandwidthOut int32  `db:"bandwidth_out"`
	InternetMaxBandwidthIn  int32  `db:"bandwidth_in"`
	SystemDiskCategory      string `db:"system_disk_category"`
	SystemDiskSize          int32  `db:"system_disk_size"`
	DataDisk                []DescribePriceRequestDataDisk
	InstanceChargeType      string `db:"instance_charge_type"`
	Renew                   int    `db:"renew"`

	SecurityGroupID string `db:"security_group_id"`
}

type GetRechargeAddressResponse struct {
	Total int
	List  []*RechargeAddress
}

type AccountRequest struct {
	Type LoginType

	Email    string
	Address  string
	Filecoin string

	Password string

	VerifyCode string

	Signature string

	InvitationCode string
}

type AccountLoginResponse struct {
	ID    string
	Token string
}

type InvitationInfo struct {
	InvCode string `db:"invitation_code"`
	ID      string `db:"id"`
}

type AccountInfo struct {
	ID         string `db:"id"`
	Email      string `db:"email"`
	Address    string `db:"address"`
	Filecoin   string `db:"filecoin"`
	Password   string `db:"passwd"`
	CreateTime int64  `db:"create_time"`
}
