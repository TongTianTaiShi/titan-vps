package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveUserInfo save user information
func (n *SQLDB) SaveUserInfo(rInfo *types.UserInfo) error {
	// update record table
	query := fmt.Sprintf(
		`INSERT INTO %s (user_addr, token) 
		        VALUES (:user_addr, :token)`, userTable)
	_, err := n.db.NamedExec(query, rInfo)

	return err
}

// UpdateUserToken update user token
func (n *SQLDB) UpdateUserToken(userAddr, token, oldToken string) error {
	query := fmt.Sprintf(`UPDATE %s SET token=? WHERE user_addr=? AND token=?`, userTable)
	_, err := n.db.Exec(query, token, userAddr, oldToken)

	return err
}

// LoadUserToken load user token
func (n *SQLDB) LoadUserToken(userAddr string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT token FROM %s WHERE user_addr=?", userTable)
	err := n.db.Get(&info, query, userAddr)
	if err != nil {
		return "0", err
	}

	return info, nil
}
