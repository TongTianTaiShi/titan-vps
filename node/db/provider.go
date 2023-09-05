package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

const accessKeyTable = "provider_access_info"

func (d *SQLDB) SaveAccessKeyInfo(info *types.AccessKeyInfo) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (provider_id, k_type, access_secret, access_key, state,  rebate, nick) 
		        VALUES (:provider_id,  :k_type, :access_secret, :access_key, :state, :rebate, :nick)`, accessKeyTable)
	_, err := d.db.NamedExec(query, info)

	return err
}

func (d *SQLDB) UpdateAccessKeyInfo(info *types.AccessKeyInfo) error {
	query := fmt.Sprintf(`UPDATE %s SET access_secret=:access_secret, access_key=:access_key, nick=:nick,
	state=:state, rebate=:rebate WHERE provider_id=:provider_id AND k_type=:k_type`, accessKeyTable)
	_, err := d.db.NamedExec(query, info)
	return err
}

func (d *SQLDB) LoadAccessKeyInfo(providerID, accessSecret string) (*types.AccessKeyInfo, error) {
	var info types.AccessKeyInfo
	query := fmt.Sprintf("SELECT * FROM %s WHERE provider_id=? AND access_secret=?", accessKeyTable)
	err := d.db.Get(&info, query, providerID, accessSecret)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// DeleteAccessKeyInfo
func (d *SQLDB) DeleteAccessKeyInfo(providerID, accessSecret string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE access_secret=? AND provider_id=?`, accessKeyTable)
	_, err := d.db.Exec(query, accessSecret, providerID)

	return err
}
