package db

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveMyInstancesInfo  save instance information
func (n *SQLDB) SaveMyInstancesInfo(rInfo *types.MyInstance) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (instance_id, order_id, user_id, private_key_status, instance_name, instance_system, location,  price,state,internet_charge_type) 
		        VALUES (:instance_id, :order_id, :user_id, :private_key_status, :instance_name, :instance_system, :location, :price,:state,:internet_charge_type)`, myInstancesTable)
	_, err := n.db.NamedExec(query, rInfo)

	return err
}

// SaveInstancesInfo save order information
func (n *SQLDB) SaveInstancesInfo(rInfo *types.DescribeInstanceTypeFromBase) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (instance_type_id, region_id, memory_size,cpu_architecture,instance_category,cpu_core_count,available_zone,instance_type_family,physical_processor_model,price,status) 
		        VALUES (:instance_type_id, :region_id, :memory_size,:cpu_architecture,:instance_category,:cpu_core_count,:available_zone,:instance_type_family,:physical_processor_model,:price,:status)
				ON DUPLICATE KEY UPDATE price=:price,status=:status`, instanceDefaultTable)
	_, err := n.db.NamedExec(query, rInfo)

	return err
}

// SaveInstancesDefaultInfo  save instance information
func (n *SQLDB) SaveInstancesDefaultInfo(rInfo *types.DescribeInstanceTypeFromBase) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (instance_type_id,region_id,price,memory_size) 
		        VALUES (:instance_type, :region_id, :price)`, instanceDefaultTable)
	_, err := n.db.NamedExec(query, rInfo)

	return err
}

// LoadMyInstancesInfo   load  my server information
func (n *SQLDB) LoadMyInstancesInfo(userID string, limit, offset int64) (*types.MyInstanceResponse, error) {
	out := new(types.MyInstanceResponse)
	var infos []*types.MyInstance
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=?  order by created_time desc LIMIT ? OFFSET ?", myInstancesTable)
	if limit > loadOrderRecordsDefaultLimit {
		limit = loadOrderRecordsDefaultLimit
	}
	err := n.db.Select(&infos, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=?", myInstancesTable)
	var count int
	err = n.db.Get(&count, countQuery, userID)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

func (n *SQLDB) LoadInstanceDetailsInfo(userID, instanceId string) (*types.InstanceDetails, error) {
	var info types.InstanceDetails
	query := fmt.Sprintf("SELECT region_id,instance_id,instance_type,image_id,security_group_id,instance_charge_type,internet_charge_type,bandwidth_out,bandwidth_in,system_disk_size,ip_address,system_disk_category,created_time,memory,memory_used,cores,cores_used,os_type FROM %s WHERE user_id=? and instance_id=?", instancesDetailsTable)
	err := n.db.Get(&info, query, userID, instanceId)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
func (n *SQLDB) LoadInstanceDefaultInfo(req *types.InstanceTypeFromBaseReq) (*types.InstanceTypeResponse, error) {
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", instanceDefaultTable, query)
	var count int
	err := n.db.Get(&count, countQuery, args...)
	if err != nil {
		return nil, err
	}
	querySql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %d OFFSET %d ", instanceDefaultTable, query, req.Limit, req.Offset)
	err = n.db.Select(&info, querySql, args...)
	if err != nil {
		return nil, err
	}
	out.Total = count
	out.List = info
	return out, nil
}
