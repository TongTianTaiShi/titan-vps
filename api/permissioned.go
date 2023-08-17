package api

import (
	"github.com/filecoin-project/go-jsonrpc/auth"
)

const (
	// When changing these, update docs/API.md too

	RoleAdmin   auth.Permission = "admin" // Manage permissions
	RoleDefault auth.Permission = "default"
	RoleUser    auth.Permission = "user"
)

var AllPermissions = []auth.Permission{RoleAdmin, RoleDefault, RoleUser}

func permissionedProxies(in, out interface{}) {
	outs := GetInternalStructs(out)
	for _, o := range outs {
		PermissionedProxy(AllPermissions, RoleDefault, in, o)
	}
}

func PermissionedTransactionAPI(a Transaction) Transaction {
	var out TransactionStruct
	permissionedProxies(a, &out)
	return &out
}

func PermissionedMallAPI(a Mall) Mall {
	var out MallStruct
	permissionedProxies(a, &out)
	return &out
}
