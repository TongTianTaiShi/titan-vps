package db

import (
	"context"
	"fmt"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/jmoiron/sqlx"
)

// LoadInstanceInfoByID loads VPS information by VPS ID.
func (d *SQLDB) LoadInstanceInfoByID(vpsID int64) (*types.InstanceDetails, error) {
	var info types.InstanceDetails
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", userInstancesTable)
	err := d.db.Get(&info, query, vpsID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// LoadUserInstanceInfoByInstanceID loads VPS information by instance ID.
func (d *SQLDB) LoadUserInstanceInfoByInstanceID(instanceID string) (*types.InstanceDetails, error) {
	var info types.InstanceDetails
	query := fmt.Sprintf("SELECT * FROM %s WHERE instance_id=?", userInstancesTable)
	err := d.db.Get(&info, query, instanceID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// SaveInstanceInfoOfUser saves VPS instance information into the database.
func (d *SQLDB) SaveInstanceInfoOfUser(rInfo *types.InstanceDetails) (int64, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (region_id,instance_id,user_id, instance_type, image_id, order_id,
			    security_group_id, instance_charge_type,internet_charge_type, period_unit, period, bandwidth_out,bandwidth_in,
			    ip_address,value,system_disk_category,system_disk_size,os_type,data_disk,auto_renew, access_key, state) 
				VALUES (:region_id,:instance_id,:user_id, :instance_type, :image_id, :order_id,
				:security_group_id, :instance_charge_type,:internet_charge_type, :period_unit, :period, :bandwidth_out,:bandwidth_in,
				:ip_address,:value,:system_disk_category,:system_disk_size,:os_type,:data_disk,:auto_renew, :access_key, :state)`, userInstancesTable)

	result, err := d.db.NamedExec(query, rInfo)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateInstanceInfoOfUser updates VPS instance information in the database.
func (d *SQLDB) UpdateInstanceInfoOfUser(info *types.InstanceDetails) error {
	query := fmt.Sprintf(`UPDATE %s SET ip_address=?, instance_id=?, os_type=?,cores=?,memory=?,expired_time=?,
	    security_group_id=?,access_key=?,bandwidth_out=?,instance_name=?,state=?,renew=?,update_time=NOW() WHERE id=?`, userInstancesTable)
	_, err := d.db.Exec(query, info.IpAddress, info.InstanceId, info.OSType, info.Cores, info.Memory, info.ExpiredTime,
		info.SecurityGroupId, info.AccessKey, info.BandwidthOut, info.InstanceName, info.State, info.Renew, info.ID)

	return err
}

// UpdateInstanceState updates VPS instance state in the database.
func (d *SQLDB) UpdateInstanceState(instanceID, state string) error {
	query := fmt.Sprintf(`UPDATE %s SET state=? WHERE instance_id=?`, userInstancesTable)
	_, err := d.db.Exec(query, state, instanceID)

	return err
}

// RenewVpsInstance updates VPS instance renewal information in the database.
func (d *SQLDB) RenewVpsInstance(info *types.InstanceDetails) error {
	query := fmt.Sprintf(`UPDATE %s SET period_unit=?, period=?, value=?,auto_renew=? WHERE instance_id=?`, userInstancesTable)
	_, err := d.db.Exec(query, info.PeriodUnit, info.Period, info.Value, info.AutoRenew, info.InstanceId)

	return err
}

// UpdateRenewInstanceStatus updates VPS instance renewal status in the database.
func (d *SQLDB) UpdateRenewInstanceStatus(info *types.SetRenewOrderReq) error {
	query := fmt.Sprintf(`UPDATE %s SET auto_renew=? WHERE instance_id=?`, userInstancesTable)
	_, err := d.db.Exec(query, info.Renew, info.InstanceId)
	if err != nil {
		return err
	}
	return nil
}

// UpdateVpsInstanceName updates VPS instance name in the database.
func (d *SQLDB) UpdateVpsInstanceName(instanceID, instanceName, userID string) error {
	query := fmt.Sprintf(`UPDATE %s SET instance_name=? WHERE instance_id=? and user_id=?`, userInstancesTable)
	_, err := d.db.Exec(query, instanceName, instanceID, userID)
	if err != nil {
		return err
	}
	//query = fmt.Sprintf(`UPDATE %s SET instance_name=? WHERE instance_id=? and user_id=?`, myInstancesTable)
	//_, err = d.db.Exec(query, instanceName, instanceID, userID)

	return err
}

// SaveInstancesInfo saves order information.
func (d *SQLDB) SaveInstancesInfo(rInfo *types.DescribeInstanceTypeFromBase) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (instance_type_id, region_id, memory_size,cpu_architecture,instance_category,cpu_core_count,available_zone,instance_type_family,physical_processor_model,price,original_price,status) 
		        VALUES (:instance_type_id, :region_id, :memory_size,:cpu_architecture,:instance_category,:cpu_core_count,:available_zone,:instance_type_family,:physical_processor_model,:price,:original_price,:status)
				ON DUPLICATE KEY UPDATE price=:price,status=:status,original_price=:original_price,updated_time=NOW()`, instanceBaseInfoTable)

	_, err := d.db.NamedExec(query, rInfo)

	return err
}

// InstancesDefaultExists checks if instance defaults exist for a specific instance type and region.
func (d *SQLDB) InstancesDefaultExists(instanceTypeID, regionID string) (bool, error) {
	var total int64
	timeString := time.Now().Format("2006-01-02")
	countSQL := fmt.Sprintf(`SELECT count(1) FROM %s WHERE instance_type_id=? and region_id=? and updated_time>?`, instanceBaseInfoTable)
	if err := d.db.Get(&total, countSQL, instanceTypeID, regionID, timeString); err != nil {
		return false, err
	}

	return total > 0, nil
}

// UpdateInstanceDefaultStatus updates the status of instance defaults for a specific instance type and region.
func (d *SQLDB) UpdateInstanceDefaultStatus(instanceTypeID, regionID string) error {
	query := fmt.Sprintf(`UPDATE %s SET status='' and updated_time=NOW() WHERE instance_type_id=? and region_id=?`, instanceBaseInfoTable)
	_, err := d.db.Exec(query, instanceTypeID, regionID)
	if err != nil {
		return err
	}

	return err
}

// LoadActiveInstancesInfo loads active instance information.
func (d *SQLDB) LoadActiveInstancesInfo(limit, page int64) (*types.GetInstanceResponse, error) {
	out := new(types.GetInstanceResponse)

	var infos []*types.InstanceDetails
	query := fmt.Sprintf("SELECT * FROM %s WHERE state!='' AND instance_id!='' order by created_time desc LIMIT ? OFFSET ?", userInstancesTable)
	if limit > loadInstancesDefaultLimit {
		limit = loadInstancesDefaultLimit
	}
	err := d.db.Select(&infos, query, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE state!='' AND instance_id!='' ", userInstancesTable)
	var count int
	err = d.db.Get(&count, countQuery)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

// // LoadInstancesInfoByUser loads user instance information.
// func (d *SQLDB) LoadInstancesInfoByUser(userID string, limit, page int64) (*types.GetInstanceResponse, error) {
// 	out := new(types.GetInstanceResponse)

// 	var infos []*types.InstanceDetails
// 	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND instance_id!='' order by created_time desc LIMIT ? OFFSET ?", userInstancesTable)
// 	if limit > loadInstancesDefaultLimit {
// 		limit = loadInstancesDefaultLimit
// 	}
// 	err := d.db.Select(&infos, query, userID, limit, page*limit)
// 	if err != nil {
// 		return nil, err
// 	}

// 	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=? AND instance_id!='' ", userInstancesTable)
// 	var count int
// 	err = d.db.Get(&count, countQuery, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	out.Total = count
// 	out.List = infos

// 	return out, nil
// }

// LoadInstancesInfoByUser loads user instance information.
func (d *SQLDB) LoadInstancesInfoByUser(userID string, limit, page int64) (*sqlx.Rows, int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=? AND instance_id!='' ", userInstancesTable)
	var total int
	err := d.db.Get(&total, countQuery, userID)
	if err != nil {
		return nil, total, err
	}

	if limit > loadInstancesDefaultLimit {
		limit = loadInstancesDefaultLimit
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND instance_id!='' order by created_time desc LIMIT ? OFFSET ?", userInstancesTable)
	rows, err := d.db.QueryxContext(context.Background(), query, userID, limit, page*limit)

	return rows, total, err
}

// LoadInstancesInfoByAccessKey loads instance information.
func (d *SQLDB) LoadInstancesInfoByAccessKey(accessKey string, limit, page int64) (*sqlx.Rows, int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE access_key=? AND instance_id!='' ", userInstancesTable)
	var total int
	err := d.db.Get(&total, countQuery, accessKey)
	if err != nil {
		return nil, total, err
	}

	if limit > loadInstancesDefaultLimit {
		limit = loadInstancesDefaultLimit
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE access_key=? AND instance_id!='' order by created_time desc LIMIT ? OFFSET ?", userInstancesTable)
	rows, err := d.db.QueryxContext(context.Background(), query, accessKey, limit, page*limit)

	return rows, total, err
}

// LoadInstancesInfo loads instance information.
func (d *SQLDB) LoadInstancesInfo(limit, page int64) (*sqlx.Rows, int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s  WHERE instance_id!=''", userInstancesTable)
	var total int
	err := d.db.Get(&total, countQuery)
	if err != nil {
		return nil, total, err
	}

	if limit > loadInstancesDefaultLimit {
		limit = loadInstancesDefaultLimit
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE instance_id!='' order by created_time desc LIMIT ? OFFSET ?", userInstancesTable)
	rows, err := d.db.QueryxContext(context.Background(), query, limit, page*limit)

	return rows, total, err
}

// DeleteInstanceInfo delete instance info by id
func (d *SQLDB) DeleteInstanceInfo(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=? `, userInstancesTable)
	_, err := d.db.Exec(query, id)

	return err
}

// LoadInstanceInfoByUser loads details of a specific instance.
func (d *SQLDB) LoadInstanceInfoByUser(userID, instanceID string) (*types.InstanceDetails, error) {
	var info types.InstanceDetails
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? and instance_id=?", userInstancesTable)
	err := d.db.Get(&info, query, userID, instanceID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// LoadInstanceDefaultInfo loads instance type defaults based on specified criteria.
func (d *SQLDB) LoadInstanceDefaultInfo(req *types.InstanceTypeFromBaseReq) (*types.InstanceTypeResponse, error) {
	out := new(types.InstanceTypeResponse)
	var info []*types.DescribeInstanceTypeFromBase
	var query string
	var args []interface{}

	query = "region_id=?"
	args = append(args, req.RegionId)
	if req.InstanceCategory != "" {
		query += " and instance_category=?"
		args = append(args, req.InstanceCategory)
	}
	if req.MemorySize != 0 {
		query += " and memory_size=?"
		args = append(args, req.MemorySize)
	}
	if req.CpuCoreCount != 0 {
		query += " and cpu_core_count=?"
		args = append(args, req.CpuCoreCount)
	}
	if req.CpuArchitecture != "" {
		query += " and cpu_architecture=?"
		args = append(args, req.CpuArchitecture)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s and status!=''", instanceBaseInfoTable, query)
	var count int
	err := d.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, err
	}

	querySQL := fmt.Sprintf("SELECT * FROM %s WHERE %s and status!='' LIMIT %d OFFSET %d ", instanceBaseInfoTable, query, req.Limit, req.Offset)
	err = d.db.Select(&info, querySQL, args...)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = info
	return out, nil
}

// LoadInstanceCPUInfo loads distinct CPU core counts based on specified criteria.
func (d *SQLDB) LoadInstanceCPUInfo(req *types.InstanceTypeFromBaseReq) ([]*int32, error) {
	var info []*int32
	var query string
	var args []interface{}
	query = "region_id=?"
	args = append(args, req.RegionId)
	if req.InstanceCategory != "" {
		query += " and instance_category=?"
		args = append(args, req.InstanceCategory)
	}
	if req.CpuCoreCount != 0 {
		query += " and cpu_core_count=?"
		args = append(args, req.CpuCoreCount)
	}
	if req.CpuArchitecture != "" {
		query += " and cpu_architecture=?"
		args = append(args, req.CpuArchitecture)
	}

	querySQL := fmt.Sprintf("SELECT distinct cpu_core_count FROM %s WHERE %s order by cpu_core_count asc", instanceBaseInfoTable, query)
	err := d.db.Select(&info, querySQL, args...)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// LoadInstanceMemoryInfo loads distinct memory sizes based on specified criteria.
func (d *SQLDB) LoadInstanceMemoryInfo(req *types.InstanceTypeFromBaseReq) ([]*float32, error) {
	var info []*float32
	var query string
	var args []interface{}
	query = "region_id=?"
	args = append(args, req.RegionId)
	if req.InstanceCategory != "" {
		query += " and instance_category=?"
		args = append(args, req.InstanceCategory)
	}
	if req.CpuCoreCount != 0 {
		query += " and cpu_core_count=?"
		args = append(args, req.CpuCoreCount)
	}
	if req.CpuArchitecture != "" {
		query += " and cpu_architecture=?"
		args = append(args, req.CpuArchitecture)
	}

	querySQL := fmt.Sprintf("SELECT distinct memory_size FROM %s WHERE %s order by memory_size asc", instanceBaseInfoTable, query)
	err := d.db.Select(&info, querySQL, args...)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// SaveInstanceRefundInfo Save administrators who unsubscribe from instances
func (d *SQLDB) SaveInstanceRefundInfo(instanceID, executor string) error {
	date := time.Now().Format("2006-01-02 15:04:05")

	query := fmt.Sprintf(
		`INSERT INTO %s (instance_id,executor,refund_time) VALUES (?,?,?)`, instanceRefundTable)
	_, err := d.db.Exec(query, instanceID, executor, date)

	return err
}

// LoadInstanceRefundInfo loads details of a specific instance.
func (d *SQLDB) LoadInstanceRefundInfo(instanceID string) (*types.InstanceDetails, error) {
	var info types.InstanceDetails
	query := fmt.Sprintf("SELECT * FROM %s WHERE instance_id=?", instanceRefundTable)
	err := d.db.Get(&info, query, instanceID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
