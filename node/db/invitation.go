package db

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
)

func (d *SQLDB) InsertInvitationCode(code string) error {
	info := types.InvitationInfo{InvCode: code}
	query := fmt.Sprintf(
		`INSERT INTO %s (invitation_code) VALUES (:invitation_code)`, invitationTable)
	_, err := d.db.NamedExec(query, info)

	return err
}

func (d *SQLDB) UpdateInvitationUUID(code, uuid string) error {
	info := types.InvitationInfo{
		InvCode: code,
		ID:      uuid,
	}

	query := fmt.Sprintf(
		`UPDATE %s SET id = :id WHERE id IS NULL AND invitation_code = :invitation_code`, invitationTable)
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
