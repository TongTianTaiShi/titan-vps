package api

import (
	"github.com/filecoin-project/go-jsonrpc/auth"
)

const (
	// When changing these, update docs/API.md too

	PermRead  auth.Permission = "read" // default
	PermWrite auth.Permission = "write"
	PermSign  auth.Permission = "sign"  // Use wallet keys for signing
	PermAdmin auth.Permission = "admin" // Manage permissions
)

var (
	AllPermissions = []auth.Permission{PermRead, PermWrite, PermSign, PermAdmin}
	DefaultPerms   = []auth.Permission{PermRead}
)

func permissionedProxies(in, out interface{}) {
	outs := GetInternalStructs(out)
	for _, o := range outs {
		auth.PermissionedProxy(AllPermissions, DefaultPerms, in, o)
	}
}

func PermissionedTransactionAPI(a Transaction) Transaction {
	var out TransactionStruct
	permissionedProxies(a, &out)
	return &out
}

func PermissionedBasisAPI(a Basis) Basis {
	var out BasisStruct
	permissionedProxies(a, &out)
	return &out
}
