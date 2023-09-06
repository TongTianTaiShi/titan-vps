package db

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
	"time"
	//"time"
)

func (d *SQLDB) SaveProviderInfo(id, userID string, loginType types.LoginType) error {
	now := time.Now().Unix()
	switch loginType {
	case types.LoginTypeMetaMask:
		return d.saveProviderInfoByAddress(id, userID, now)
	case types.LoginTypeEmail:
		return d.saveProviderInfoByEmail(id, userID, now)
	case types.LoginTypeFilecoin:
		return d.saveProviderInfoByFilecoin(id, userID, now)
	default:
		return fmt.Errorf("%s not  supported", loginType.String())
	}
}

func (d *SQLDB) saveProviderInfoByEmail(id, email string, now int64) error {
	info := types.AccountInfo{ID: id, Email: email, CreateTime: now}
	query := fmt.Sprintf(
		`INSERT INTO %s (id, email, create_time) VALUES (:id, :email, :create_time)`, providerInfoTable)
	_, err := d.db.NamedExec(query, info)
	return err
}

func (d *SQLDB) saveProviderInfoByAddress(id, address string, now int64) error {
	info := types.AccountInfo{ID: id, Address: address, CreateTime: now}
	query := fmt.Sprintf(
		`INSERT INTO %s (id, address, create_time) VALUES (:id, :address, :create_time)`, providerInfoTable)
	_, err := d.db.NamedExec(query, info)
	return err
}

func (d *SQLDB) saveProviderInfoByFilecoin(id, address string, now int64) error {
	info := types.AccountInfo{ID: id, Address: address, CreateTime: now}
	query := fmt.Sprintf(
		`INSERT INTO %s (id, filecoin, create_time) VALUES (:id, :filecoin, :create_time)`, providerInfoTable)
	_, err := d.db.NamedExec(query, info)
	return err
}

func (d *SQLDB) CheckMerchantIsExist(userID string, loginType types.LoginType) (bool, error) {
	switch loginType {
	case types.LoginTypeMetaMask:
		return d.checkAddressIsExist(userID)
	case types.LoginTypeEmail:
		return d.checkEmailIsExist(userID)
	case types.LoginTypeFilecoin:
		return d.checkFilecoinIsExist(userID)
	default:
		return false, fmt.Errorf("%s not  supported", loginType.String())
	}
}

func (d *SQLDB) checkEmailIsExist(email string) (bool, error) {
	var count int
	err := d.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE email = ?`, providerInfoTable), email).Scan(&count)

	return checkIsExist(count, err)
}

func (d *SQLDB) checkAddressIsExist(address string) (bool, error) {
	var count int
	err := d.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE address = ?`, providerInfoTable), address).Scan(&count)

	return checkIsExist(count, err)
}

func (d *SQLDB) checkFilecoinIsExist(address string) (bool, error) {
	var count int
	err := d.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE filecoin = ?`, providerInfoTable), address).Scan(&count)

	return checkIsExist(count, err)
}

func checkIsExist(count int, err error) (bool, error) {
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *SQLDB) GetProviderUUID(userID string, loginType types.LoginType) (string, error) {
	switch loginType {
	case types.LoginTypeMetaMask:
		return d.getProviderInfoByAddress(userID)
	case types.LoginTypeEmail:
		return d.getProviderInfoByEmail(userID)
	case types.LoginTypeFilecoin:
		return d.getProviderInfoByFilecoin(userID)
	default:
		return "", fmt.Errorf("%s not  supported", loginType.String())
	}
}

func (d *SQLDB) getProviderInfoByAddress(address string) (string, error) {
	var id string
	query := fmt.Sprintf(
		`SELECT id FROM %s WHERE address = ?`, providerInfoTable)
	err := d.db.QueryRow(query, address).Scan(&id)
	return id, err
}

func (d *SQLDB) getProviderInfoByEmail(email string) (string, error) {
	var id string
	query := fmt.Sprintf(
		`SELECT id FROM %s WHERE email = ?`, providerInfoTable)
	err := d.db.QueryRow(query, email).Scan(&id)
	return id, err
}

func (d *SQLDB) getProviderInfoByFilecoin(address string) (string, error) {
	var id string
	query := fmt.Sprintf(
		`SELECT id FROM %s WHERE filecoin = ?`, providerInfoTable)
	err := d.db.QueryRow(query, address).Scan(&id)
	return id, err
}
