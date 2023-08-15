package db

import (
	"fmt"
)

// SaveRechargeAddress save user information
func (n *SQLDB) SaveRechargeAddress(addresses []string) error {
	// update record table
	query := fmt.Sprintf(
		`INSERT INTO %s (addr) VALUES (?)`, rechargeAddressTable)
	_, err := n.db.NamedExec(query, addresses)

	return err
}

// UpdateRechargeAddressOfUser save user information
func (n *SQLDB) UpdateRechargeAddressOfUser(addr, userAddr string) error {
	// update record table
	dQuery := fmt.Sprintf(`UPDATE %s SET user_addr=? WHERE addr=? AND user_addr="" `, rechargeAddressTable)
	_, err := n.db.Exec(dQuery, userAddr, addr)

	return err
}

// GetUserOfRechargeAddress get user address
func (n *SQLDB) GetUserOfRechargeAddress(addr string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT user_addr FROM %s WHERE addr=?", rechargeAddressTable)
	err := n.db.Get(&info, query, addr)
	if err != nil {
		return "", err
	}

	return info, nil
}

// GetRechargeAddressOfUser get user recharge address
func (n *SQLDB) GetRechargeAddressOfUser(userAddr string) (string, error) {
	var info string
	query := fmt.Sprintf("SELECT addr FROM %s WHERE user_addr=?", rechargeAddressTable)
	err := n.db.Get(&info, query, userAddr)
	if err != nil {
		return "", err
	}

	return info, nil
}
