package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// LoadVpsInfo  load  vps information
func (n *SQLDB) LoadVpsInfo(vpsID int64) (*types.CreateInstanceReq, error) {
	var info types.CreateInstanceReq
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", vpsInstanceTable)
	err := n.db.Get(&info, query, vpsID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (n *SQLDB) LoadVpsDeviceInfo(instanceID int64) (*types.CreateInstanceReq, error) {
	var info types.CreateInstanceReq
	query := fmt.Sprintf("SELECT * FROM %s WHERE instance_id=?", vpsInstanceDeviceTable)
	err := n.db.Get(&info, query, instanceID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// VpsExists  checks if this vps info exists in the state machine table of the specified server.
func (n *SQLDB) VpsExists(vpsID int64) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(id) FROM %s WHERE id=? `, vpsInstanceTable)
	if err := n.db.Get(&total, countSQL, vpsID); err != nil {
		return false, err
	}

	return total > 0, nil
}

func (n *SQLDB) VpsDeviceExists(instanceID int64) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(instance_id) FROM %s WHERE instance_id=? `, vpsInstanceTable)
	if err := n.db.Get(&total, countSQL, instanceID); err != nil {
		return false, err
	}

	return total > 0, nil
}

// SaveVpsInstance   saves vps info into the database.
func (n *SQLDB) SaveVpsInstance(rInfo *types.CreateInstanceReq) (int64, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (region_id, instance_type, dry_run, image_id, security_group_id, instanceCharge_type, period_unit, period, bandwidth_out,bandwidth_in) 
				VALUES (:region_id, :instance_type, :dry_run, :image_id, :security_group_id, :instanceCharge_type, :period_unit, :period, :bandwidth_out,bandwidth_in)`, vpsInstanceTable)

	result, err := n.db.NamedExec(query, rInfo)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (n *SQLDB) SaveVpsInstanceDevice(rInfo *types.CreateInstanceResponse) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (instance_id, order_id, request_id, trade_price, public_ip_address) 
					VALUES (:instance_id, :order_id, :request_id, :trade_price, :public_ip_address)`, vpsInstanceTable)

	_, err := n.db.NamedExec(query, rInfo)
	if err != nil {
		return err
	}

	return nil
}
