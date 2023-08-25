package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// LoadVpsInfo  load  vps information
func (n *SQLDB) LoadVpsInfo(vpsID int64) (*types.CreateInstanceReq, error) {
	var info types.CreateInstanceReq
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", instancesDetailsTable)
	err := n.db.Get(&info, query, vpsID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (n *SQLDB) LoadVpsInfoByInstanceId(instanceID string) (*types.CreateInstanceReq, error) {
	var info types.CreateInstanceReq
	query := fmt.Sprintf("SELECT * FROM %s WHERE instance_id=?", instancesDetailsTable)
	err := n.db.Get(&info, query, instanceID)
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
	countSQL := fmt.Sprintf(`SELECT count(id) FROM %s WHERE id=? `, instancesDetailsTable)
	if err := n.db.Get(&total, countSQL, vpsID); err != nil {
		return false, err
	}

	return total > 0, nil
}

func (n *SQLDB) VpsDeviceExists(instanceID int64) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(instance_id) FROM %s WHERE instance_id=? `, instancesDetailsTable)
	if err := n.db.Get(&total, countSQL, instanceID); err != nil {
		return false, err
	}

	return total > 0, nil
}

// SaveVpsInstance   saves vps info into the database.
func (n *SQLDB) SaveVpsInstance(rInfo *types.CreateOrderReq) (int64, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (region_id,instance_id,user_id,order_id, instance_type, dry_run, image_id, security_group_id, instance_charge_type,internet_charge_type, period_unit, period, bandwidth_out,bandwidth_in,ip_address,trade_price,system_disk_category,system_disk_size,os_type,data_disk,renew) 
				VALUES (:region_id,:instance_id,:user_id,:order_id, :instance_type, :dry_run, :image_id, :security_group_id, :instance_charge_type,:internet_charge_type, :period_unit, :period, :bandwidth_out,:bandwidth_in,:ip_address,:trade_price,:system_disk_category,:system_disk_size,:os_type,:data_disk,:renew)`, instancesDetailsTable)

	result, err := n.db.NamedExec(query, rInfo)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (n *SQLDB) UpdateVpsInstance(info *types.CreateInstanceReq) error {
	query := fmt.Sprintf(`UPDATE %s SET ip_address=?, instance_id=?, user_id=?,os_type=?,cores=?,memory=?,
	    security_group_id=? WHERE order_id=?`, instancesDetailsTable)
	_, err := n.db.Exec(query, info.IpAddress, info.InstanceId, info.UserID, info.OSType, info.Cores, info.Memory, info.SecurityGroupId, info.OrderID)

	return err
}

func (n *SQLDB) RenewVpsInstance(info *types.CreateInstanceReq) error {
	query := fmt.Sprintf(`UPDATE %s SET period_unit=?, period=?, trade_price=?,renew=? WHERE instance_id=?`, instancesDetailsTable)
	_, err := n.db.Exec(query, info.PeriodUnit, info.Period, info.TradePrice, info.Renew, info.InstanceId)
	if err != nil {
		return err
	}
	return nil
}
func (n *SQLDB) UpdateRenewInstanceStatus(info *types.SetRenewOrderReq) error {
	query := fmt.Sprintf(`UPDATE %s SET renew=? WHERE instance_id=?`, instancesDetailsTable)
	_, err := n.db.Exec(query, info.Renew, info.InstanceId)
	if err != nil {
		return err
	}
	return nil
}
func (n *SQLDB) UpdateVpsInstanceName(instanceID, instanceName, userID string) error {
	query := fmt.Sprintf(`UPDATE %s SET instance_name=? WHERE instance_id=? and user_id=?`, instancesDetailsTable)
	_, err := n.db.Exec(query, instanceName, instanceID, userID)
	if err != nil {
		return err
	}
	query = fmt.Sprintf(`UPDATE %s SET instance_name=? WHERE instance_id=? and user_id=?`, myInstancesTable)
	_, err = n.db.Exec(query, instanceName, instanceID, userID)

	return err
}

func (n *SQLDB) SaveVpsInstanceDevice(rInfo *types.CreateInstanceResponse) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (instance_id, order_id, request_id, trade_price, public_ip_address) 
					VALUES (:instance_id, :order_id, :request_id, :trade_price, :public_ip_address)`, vpsInstanceDeviceTable)

	_, err := n.db.NamedExec(query, rInfo)
	if err != nil {
		return err
	}

	return nil
}
