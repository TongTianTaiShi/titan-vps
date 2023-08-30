package db

import (
	"context"
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/jmoiron/sqlx"
)

// SaveOrderInfo saves order information.
func (d *SQLDB) SaveOrderInfo(rInfo *types.OrderRecord) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (order_id, value, state, done_state, vps_id, msg, user_id, order_type, end_time) 
		        VALUES (:order_id, :value, :state, :done_state, :vps_id, :msg, :user_id, :order_type, :end_time)
				ON DUPLICATE KEY UPDATE state=:state, done_state=:done_state, done_time=NOW(), user_id=:user_id,
				value=:value, vps_id=:vps_id, msg=:msg, order_type=:order_type`, orderRecordTable)
	_, err := d.db.NamedExec(query, rInfo)

	return err
}

// LoadOrderRecord loads order record information.
func (d *SQLDB) LoadOrderRecord(orderID string, minute int64) (*types.OrderRecord, error) {
	var info types.OrderRecord
	query := fmt.Sprintf("SELECT *,DATE_ADD(created_time,INTERVAL %d MINUTE) AS expiration FROM %s WHERE order_id=?", minute, orderRecordTable)
	err := d.db.Get(&info, query, orderID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// LoadOrderRecordByUserUndone loads undone order records for a specific user with pagination.
func (d *SQLDB) LoadOrderRecordByUserUndone(userID string, limit, page int64, minute int) (*types.OrderRecordResponse, error) {
	out := new(types.OrderRecordResponse)
	var infos []*types.OrderRecord
	query := fmt.Sprintf("SELECT *,DATE_ADD(created_time,INTERVAL %d MINUTE) AS expiration FROM %s WHERE user_id=? and state!=3  order by created_time desc LIMIT ? OFFSET ?", minute, orderRecordTable)
	if limit > loadOrderRecordsDefaultLimit {
		limit = loadOrderRecordsDefaultLimit
	}
	err := d.db.Select(&infos, query, userID, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=? and state!=3", orderRecordTable)
	var count int
	err = d.db.Get(&count, countQuery, userID)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

// LoadOrderRecordsByUser loads order records for a specific user with pagination.
func (d *SQLDB) LoadOrderRecordsByUser(userID string, limit, page int64, minute int) (*types.OrderRecordResponse, error) {
	out := new(types.OrderRecordResponse)
	var infos []*types.OrderRecord
	query := fmt.Sprintf("SELECT *,DATE_ADD(created_time,INTERVAL %d MINUTE) AS expiration FROM %s WHERE user_id=?  order by created_time desc LIMIT ? OFFSET ?", minute, orderRecordTable)
	if limit > loadOrderRecordsDefaultLimit {
		limit = loadOrderRecordsDefaultLimit
	}
	err := d.db.Select(&infos, query, userID, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=?", orderRecordTable)
	var count int
	err = d.db.Get(&count, countQuery, userID)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

// OrderExists checks if an order exists.
func (d *SQLDB) OrderExists(orderID string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(order_id) FROM %s WHERE order_id=? `, orderRecordTable)
	if err := d.db.Get(&total, countSQL, orderID); err != nil {
		return false, err
	}

	return total > 0, nil
}

// LoadOrderRecords loads order records with specified statuses, pagination, and minute interval.
func (d *SQLDB) LoadOrderRecords(statuses []int64, limit, page, minute int) (*sqlx.Rows, error) {
	if limit > loadOrderRecordsDefaultLimit {
		limit = loadOrderRecordsDefaultLimit
	}
	sQuery := fmt.Sprintf(`SELECT *,DATE_ADD(created_time,INTERVAL %d MINUTE) AS expiration FROM %s WHERE state in (?) order by order_id asc LIMIT ? OFFSET ?`, minute, orderRecordTable)
	query, args, err := sqlx.In(sQuery, statuses, limit, page*limit)
	if err != nil {
		return nil, err
	}

	query = d.db.Rebind(query)
	return d.db.QueryxContext(context.Background(), query, args...)
}

// LoadOrderCount counts the number of orders.
func (d *SQLDB) LoadOrderCount() (int, error) {
	var size int
	cmd := fmt.Sprintf("SELECT count(order_id) FROM %s", orderRecordTable)
	err := d.db.Get(&size, cmd)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// LoadAllOrderRecords loads all order records with specified statuses and minute interval.
func (d *SQLDB) LoadAllOrderRecords(statuses []int64, minute int64) (*sqlx.Rows, error) {
	sQuery := fmt.Sprintf(`SELECT *,DATE_ADD(created_time,INTERVAL %d MINUTE) AS expiration FROM %s WHERE state in (?) `, minute, orderRecordTable)
	query, args, err := sqlx.In(sQuery, statuses)
	if err != nil {
		return nil, err
	}

	query = d.db.Rebind(query)
	return d.db.QueryxContext(context.Background(), query, args...)
}
