package db

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
	"time"
	//"time"
)

func (d *SQLDB) SaveMerchantInfo(userID string, loginType types.LoginType) error {
	now := time.Now().Unix()
	switch loginType {
	case types.LoginTypeMetaMask:
		return d.saveMerchantInfoByAddress(userID, now)
	case types.LoginTypeEmail:
		return d.saveMerchantInfoByEmail(userID, now)
	case types.LoginTypeFilecoin:
		return d.saveMerchantInfoByFilecoin(userID, now)
	default:
		return fmt.Errorf("%s not  supported", loginType.String())
	}
}

func (d *SQLDB) saveMerchantInfoByEmail(email string, now int64) error {
	info := types.AccountInfo{Email: email, CreateTime: now}
	query := fmt.Sprintf(
		`INSERT INTO %s (email, create_time) VALUES (:email, :create_time)`, accountTable)
	_, err := d.db.NamedExec(query, info)
	return err
}

func (d *SQLDB) saveMerchantInfoByAddress(address string, now int64) error {
	info := types.AccountInfo{Address: address, CreateTime: now}
	query := fmt.Sprintf(
		`INSERT INTO %s (address, create_time) VALUES (:address, :create_time)`, accountTable)
	_, err := d.db.NamedExec(query, info)
	return err
}

func (d *SQLDB) saveMerchantInfoByFilecoin(address string, now int64) error {
	info := types.AccountInfo{Filecoin: address, CreateTime: now}
	query := fmt.Sprintf(
		`INSERT INTO %s (filecoin, create_time) VALUES (:filecoin, :create_time)`, accountTable)
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
		`SELECT COUNT(*) FROM %s WHERE email = ?`, accountTable), email).Scan(&count)

	return checkIsExist(count, err)
}

func (d *SQLDB) checkAddressIsExist(address string) (bool, error) {
	var count int
	err := d.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE address = ?`, accountTable), address).Scan(&count)

	return checkIsExist(count, err)
}

func (d *SQLDB) checkFilecoinIsExist(address string) (bool, error) {
	var count int
	err := d.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE filecoin = ?`, accountTable), address).Scan(&count)

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
