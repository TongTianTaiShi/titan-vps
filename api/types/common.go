package types

import (
	"github.com/filecoin-project/go-jsonrpc/auth"
)

type OpenRPCDocument map[string]interface{}

type JWTPayload struct {
	Allow []auth.Permission
	ID    string

	LoginType int64

	// Extend is json string
	Extend string
}
