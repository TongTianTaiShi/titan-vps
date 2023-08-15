package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveUserInfo save user information
func (n *SQLDB) SaveUserInfo(rInfo *types.UserInfo) error {
	// update record table
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, balance) 
		        VALUES (:user_id, :balance)`, userTable)
	_, err := n.db.NamedExec(query, rInfo)

	return err
}

// UpdateUserBalance update user balance
func (n *SQLDB) UpdateUserBalance(userID, balance, oldBalance string) error {
	query := fmt.Sprintf(`UPDATE %s SET balance=? WHERE user_id=? AND balance=?`, userTable)
	_, err := n.db.Exec(query, balance, userID, oldBalance)

	return err
}

// LoadUserBalance load user balance
func (n *SQLDB) LoadUserBalance(userID string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=?", userTable)
	err := n.db.Get(&info, query, userID)
	if err != nil {
		return "0", err
	}

	return info, nil
}
