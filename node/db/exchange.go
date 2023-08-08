package db

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
)

// SaveRechargeInfo save recharge information
func (n *SQLDB) SaveRechargeInfo(rInfo *types.RechargeRecord) error {
	// update record table
	query := fmt.Sprintf(
		`INSERT INTO %s (id, from_addr, to_addr, value, created_height, done_height, state, recharge_addr, recharge_hash, msg, user_addr, tx_hash, done_state) 
		        VALUES (:id, :from_addr, :to_addr, :value, :created_height, :done_height, :state, :recharge_addr, :recharge_hash, :msg, :user_addr, :tx_hash, :done_state)`, rechargeRecordTable)
	_, err := n.db.NamedExec(query, rInfo)

	return err
}

// LoadRechargeRecord load recharge record information
func (n *SQLDB) LoadRechargeRecord(id string) (*types.RechargeRecord, error) {
	var info types.RechargeRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", rechargeRecordTable)
	err := n.db.Get(&info, query, id)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// UpdateRechargeRecord update recharge record information
func (n *SQLDB) UpdateRechargeRecord(info *types.RechargeRecord, newState types.RechargeState) error {
	query := fmt.Sprintf(`UPDATE %s SET state=?, value=?, done_state=?, done_time=NOW(), from_addr=?,
	    done_height=?, tx_hash=?, recharge_hash=?, msg=? WHERE id=? AND state=?`, rechargeRecordTable)
	_, err := n.db.Exec(query, newState, info.Value, info.DoneState, info.From, info.DoneHeight, info.TxHash, info.RechargeHash, info.Msg, info.ID, info.State)

	return err
}

// LoadRechargeRecords load the recharge records from the incoming scheduler
func (n *SQLDB) LoadRechargeRecords(state types.RechargeState) ([]*types.RechargeRecord, error) {
	var infos []*types.RechargeRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE state=? ", rechargeRecordTable)

	err := n.db.Select(&infos, query, state)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
