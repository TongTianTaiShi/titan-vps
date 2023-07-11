package api

import (
	"context"
)

// Transaction is an interface for transaction
type Transaction interface {
	Common

	Hello(ctx context.Context) error //perm:read
}
