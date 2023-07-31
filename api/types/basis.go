package types

import "time"

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
