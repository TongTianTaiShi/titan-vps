package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/jmoiron/sqlx"
)

// UpdateOrderInfo update asset information
func (n *SQLDB) UpdateOrderInfo(hash string, state, doneState, doneHeight int64) error {
	tx, err := n.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Errorf("SaveAssetRecord Rollback err:%s", err.Error())
		}
	}()

	fmt.Println("UpdateAssetInfo state : ", state)

	// update record table
	dQuery := fmt.Sprintf(`UPDATE %s SET state=?, done_height=?, done_state=?, done_time=NOW() WHERE hash=?`, orderRecordTable)
	_, err = tx.Exec(dQuery, state, doneHeight, doneState, hash)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// LoadOrderRecord load asset record information
func (n *SQLDB) LoadOrderRecord(hash string) (*types.OrderRecord, error) {
	var info types.OrderRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE hash=?", orderRecordTable)
	err := n.db.Get(&info, query, hash)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// AssetExists checks if an asset exists in the state machine table of the specified server.
func (n *SQLDB) AssetExists(hash string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(hash) FROM %s WHERE hash=? `, orderRecordTable)
	if err := n.db.Get(&total, countSQL, hash); err != nil {
		return false, err
	}

	return total > 0, nil
}

// LoadAssetRecords load the asset records from the incoming scheduler
func (n *SQLDB) LoadAssetRecords(statuses []int64, limit, offset int) (*sqlx.Rows, error) {
	if limit > loadAssetRecordsDefaultLimit || limit == 0 {
		limit = loadAssetRecordsDefaultLimit
	}
	sQuery := fmt.Sprintf(`SELECT * FROM %s WHERE state in (?) order by hash asc LIMIT ? OFFSET ?`, orderRecordTable)
	query, args, err := sqlx.In(sQuery, statuses, limit, offset)
	if err != nil {
		return nil, err
	}

	query = n.db.Rebind(query)
	return n.db.QueryxContext(context.Background(), query, args...)
}

// LoadAssetCount count asset
func (n *SQLDB) LoadAssetCount() (int, error) {
	var size int
	cmd := fmt.Sprintf("SELECT count(hash) FROM %s", orderRecordTable)
	err := n.db.Get(&size, cmd)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// LoadAllAssetRecords loads all asset records for a given server ID.
func (n *SQLDB) LoadAllAssetRecords(limit, offset int, statuses []int64) (*sqlx.Rows, error) {
	sQuery := fmt.Sprintf(`SELECT * FROM %s WHERE state in (?) order by hash asc limit ? offset ?`, orderRecordTable)
	query, args, err := sqlx.In(sQuery, statuses, limit, offset)
	if err != nil {
		return nil, err
	}

	query = n.db.Rebind(query)
	return n.db.QueryxContext(context.Background(), query, args...)
}

// SaveOrderRecord  saves an asset record into the database.
func (n *SQLDB) SaveOrderRecord(rInfo *types.OrderRecord) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (hash, from_addr, to_addr, value, created_height, done_height, state, done_state, vps_id) 
				VALUES (:hash, :from_addr, :to_addr, :value, :created_height, :done_height, :state, :done_state, :vps_id)`, orderRecordTable)

	_, err := n.db.NamedExec(query, rInfo)

	return err
}
