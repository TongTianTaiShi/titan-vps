package db

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/api/types"
)

// LoadMyInstancesInfo   load  my server information
func (n *SQLDB) LoadMyInstancesInfo(userId string) (*types.MyInstances, error) {
	var info types.MyInstances
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=?", myInstancesTable)
	err := n.db.Get(&info, query, userId)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (n *SQLDB) LoadInstanceDetailsInfo(instanceId string) (*types.InstanceDetails, error) {
	var info types.InstanceDetails
	query := fmt.Sprintf("SELECT * FROM %s WHERE instance_id=?", instancesDetailsTable)
	err := n.db.Get(&info, query, instanceId)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
