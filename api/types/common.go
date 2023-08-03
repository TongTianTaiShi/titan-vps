package types

import (
	"github.com/filecoin-project/go-jsonrpc/auth"
)

type OpenRPCDocument map[string]interface{}

type JWTPayload struct {
	Allow []auth.Permission
	ID    string
	// TODO remove NodeID later, any role id replace as ID
	NodeID string
	// Extend is json string
	Extend string
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
