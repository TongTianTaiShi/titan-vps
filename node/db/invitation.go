package db

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
)

func (d *SQLDB) InsertInvitationCode(code string) error {
	info := types.InvitationInfo{ID: code}
	query := fmt.Sprintf(
		`INSERT IGNORE INTO %s (id) VALUES (:id)`, invitationTable)
	_, err := d.db.NamedExec(query, info)

	return err
}

func (d *SQLDB) UpdateInvitationUserID(code, userID string) error {
	info := types.InvitationInfo{
		ID:     code,
		UserID: userID,
	}

	query := fmt.Sprintf(
		`UPDATE %s SET user_id = COALESCE(:user_id, user_id) WHERE id = :id`, invitationTable)
	result, err := d.db.NamedExec(query, info)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("undefined invitation code: %s", code)
	}

	return nil
}
