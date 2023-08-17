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
