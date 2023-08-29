package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveUserInfo saves user information.
func (d *SQLDB) SaveUserInfo(rInfo *types.UserInfo) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, balance) 
		        VALUES (:user_id, :balance)`, userTable)
	_, err := d.db.NamedExec(query, rInfo)

	return err
}

// UpdateUserBalance updates user balance.
func (d *SQLDB) UpdateUserBalance(userID, balance, oldBalance string) error {
	query := fmt.Sprintf(`UPDATE %s SET balance=? WHERE user_id=? AND balance=?`, userTable)
	_, err := d.db.Exec(query, balance, userID, oldBalance)

	return err
}

// LoadUserBalance loads user balance.
func (d *SQLDB) LoadUserBalance(userID string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=?", userTable)
	err := d.db.Get(&info, query, userID)
	if err != nil {
		return "0", err
	}

	return info, nil
}

// UserExists checks if a user exists.
func (d *SQLDB) UserExists(userID string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(user_id) FROM %s WHERE user_id=? `, userTable)
	if err := d.db.Get(&total, countSQL, userID); err != nil {
		return false, err
	}

	return total > 0, nil
}

// SaveAdminInfo saves admin information.
func (d *SQLDB) SaveAdminInfo(userID, nickName string) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, nick_name) 
		        VALUES (?, ?)`, adminTable)
	_, err := d.db.Exec(query, userID, nickName)

	return err
}

// AdminExists checks if an admin exists.
func (d *SQLDB) AdminExists(userID string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(user_id) FROM %s WHERE user_id=? `, adminTable)
	if err := d.db.Get(&total, countSQL, userID); err != nil {
		return false, err
	}

	return total > 0, nil
}
