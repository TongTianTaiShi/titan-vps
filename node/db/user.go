package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveUserInfo save user information
func (n *SQLDB) SaveUserInfo(rInfo *types.UserInfo) error {
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

// UserExists checks if an user exists
func (n *SQLDB) UserExists(userID string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(user_id) FROM %s WHERE user_id=? `, userTable)
	if err := n.db.Get(&total, countSQL, userID); err != nil {
		return false, err
	}

	return total > 0, nil
}

// SaveAdminInfo save admin information
func (n *SQLDB) SaveAdminInfo(userID, nickName string) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, nick_name) 
		        VALUES (?, ?)`, adminTable)
	_, err := n.db.Exec(query, userID, nickName)

	return err
}

// AdminExists checks if an admin exists
func (n *SQLDB) AdminExists(userID string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(user_id) FROM %s WHERE user_id=? `, adminTable)
	if err := n.db.Get(&total, countSQL, userID); err != nil {
		return false, err
	}

	return total > 0, nil
}
