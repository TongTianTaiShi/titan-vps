package db

import (
	"database/sql"
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveRechargeAddresses inserts multiple addresses into the database.
func (d *SQLDB) SaveRechargeAddresses(addresses []string) error {
	// update record table
	for _, addr := range addresses {
		query := fmt.Sprintf(
			`INSERT INTO %s (addr) VALUES (?)`, rechargeAddressTable)
		d.db.Exec(query, addr)
	}

	return nil
}

// AssignUserToRechargeAddress assigns a user to a recharge address.
func (d *SQLDB) AssignUserToRechargeAddress(addr, userID string) error {
	// update record table
	dQuery := fmt.Sprintf(`UPDATE %s SET user_id=? WHERE addr=? AND user_id="" `, rechargeAddressTable)
	_, err := d.db.Exec(dQuery, userID, addr)

	return err
}

// LoadUserByRechargeAddress retrieves the user associated with a recharge address.
func (d *SQLDB) LoadUserByRechargeAddress(addr string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE addr=?", rechargeAddressTable)
	err := d.db.Get(&info, query, addr)
	if err != nil {
		return "", err
	}

	return info, nil
}

// LoadRechargeAddressByUser retrieves the recharge address associated with a user.
func (d *SQLDB) LoadRechargeAddressByUser(userID string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT addr FROM %s WHERE user_id=?", rechargeAddressTable)
	err := d.db.Get(&info, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return info, nil
}

// LoadUnusedRechargeAddress retrieves an unused recharge address.
func (d *SQLDB) LoadUnusedRechargeAddress() (string, error) {
	var addr string
	query := fmt.Sprintf("SELECT addr FROM %s WHERE user_id='' limit 1 ", rechargeAddressTable)
	err := d.db.Get(&addr, query)
	if err != nil {
		return "", err
	}

	return addr, nil
}

// LoadUsedRechargeAddresses retrieves all user-assigned recharge addresses.
func (d *SQLDB) LoadUsedRechargeAddresses() ([]types.RechargeAddress, error) {
	var infos []types.RechargeAddress
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id !='' ", rechargeAddressTable)
	err := d.db.Select(&infos, query)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

// LoadRechargeAddresses retrieves a list of recharge addresses with pagination.
func (d *SQLDB) LoadRechargeAddresses(limit, page int64) (*types.GetRechargeAddressResponse, error) {
	out := new(types.GetRechargeAddressResponse)

	var infos []*types.RechargeAddress
	query := fmt.Sprintf("SELECT * FROM %s order by user_id desc LIMIT ? OFFSET ?", rechargeAddressTable)
	if limit > loadAddressesDefaultLimit {
		limit = loadAddressesDefaultLimit
	}

	err := d.db.Select(&infos, query, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s ", rechargeAddressTable)
	var count int
	err = d.db.Get(&count, countQuery)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}
