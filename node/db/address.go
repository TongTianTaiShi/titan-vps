package db

import (
	"database/sql"
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveRechargeAddress save user information
func (n *SQLDB) SaveRechargeAddress(addresses []string) error {
	// update record table
	for _, addr := range addresses {
		query := fmt.Sprintf(
			`INSERT INTO %s (addr) VALUES (?)`, rechargeAddressTable)
		n.db.Exec(query, addr)
	}

	return nil
}

// UpdateRechargeAddressOfUser save user information
func (n *SQLDB) UpdateRechargeAddressOfUser(addr, userID string) error {
	// update record table
	dQuery := fmt.Sprintf(`UPDATE %s SET user_id=? WHERE addr=? AND user_id="" `, rechargeAddressTable)
	_, err := n.db.Exec(dQuery, userID, addr)

	return err
}

// LoadUserOfRechargeAddress get user address
func (n *SQLDB) LoadUserOfRechargeAddress(addr string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE addr=?", rechargeAddressTable)
	err := n.db.Get(&info, query, addr)
	if err != nil {
		return "", err
	}

	return info, nil
}

// LoadRechargeAddressOfUser get user recharge address
func (n *SQLDB) LoadRechargeAddressOfUser(userID string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT addr FROM %s WHERE user_id=?", rechargeAddressTable)
	err := n.db.Get(&info, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return info, nil
}

// LoadUnusedRechargeAddress get an unused address
func (n *SQLDB) LoadUnusedRechargeAddress() (string, error) {
	var addr string
	query := fmt.Sprintf("SELECT addr FROM %s WHERE user_id='' limit 1 ", rechargeAddressTable)
	err := n.db.Get(&addr, query)
	if err != nil {
		return "", err
	}

	return addr, nil
}

// LoadUsedRechargeAddresses get user recharge address
func (n *SQLDB) LoadUsedRechargeAddresses() ([]types.RechargeAddress, error) {
	var infos []types.RechargeAddress
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id !='' ", rechargeAddressTable)
	err := n.db.Select(&infos, query)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (n *SQLDB) LoadRechargeAddresses(limit, page int64) (*types.GetRechargeAddressResponse, error) {
	out := new(types.GetRechargeAddressResponse)

	var infos []*types.RechargeAddress
	query := fmt.Sprintf("SELECT * FROM %s order by user_id desc LIMIT ? OFFSET ?", rechargeAddressTable)
	if limit > loadAddressesDefaultLimit {
		limit = loadAddressesDefaultLimit
	}

	err := n.db.Select(&infos, query, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s ", rechargeAddressTable)
	var count int
	err = n.db.Get(&count, countQuery)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}
