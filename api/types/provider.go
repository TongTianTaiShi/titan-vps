package types

import "time"

// AccessKeyType
type AccessKeyType int64

const (
	// AccessKeyAliyun
	AccessKeyAliyun AccessKeyType = iota
)

// AccessKeyState
type AccessKeyState int64

const (
	AccessKeyStateNormal AccessKeyState = iota
	AccessKeyStateExceptional
)

// AccessKeyInfo
type AccessKeyInfo struct {
	ProviderID   string         `db:"provider_id"`
	AccessSecret string         `db:"access_secret"`
	AccessKey    string         `db:"access_key"`
	Type         AccessKeyType  `db:"k_type"`
	State        AccessKeyState `db:"state"`
	CreatedTime  time.Time      `db:"created_time"`
	Rebate       float64        `db:"rebate"`
	NickName     string         `db:"nick"`
}
