package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/jmoiron/sqlx"
)

// SaveRechargeRecordAndUserBalance saves recharge information and updates user balance.
func (d *SQLDB) SaveRechargeRecordAndUserBalance(rInfo *types.RechargeRecord, balance, oldBalance string) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Errorf("SaveRechargeRecordAndUserBalance Rollback err:%s", err.Error())
		}
	}()

	query := fmt.Sprintf(
		`INSERT INTO %s (order_id, from_addr, to_addr, value, state,  user_id) 
		        VALUES (:order_id, :from_addr, :to_addr, :value, :state, :user_id)`, rechargeRecordTable)
	_, err = tx.NamedExec(query, rInfo)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`UPDATE %s SET balance=? WHERE user_id=? AND balance=?`, userTable)
	_, err = tx.Exec(query, balance, rInfo.UserID, oldBalance)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RechargeRecordExists checks if a recharge order exists.
func (d *SQLDB) RechargeRecordExists(orderID string) (bool, error) {
	var total int64
	countSQL := fmt.Sprintf(`SELECT count(order_id) FROM %s WHERE order_id=? `, rechargeRecordTable)
	if err := d.db.Get(&total, countSQL, orderID); err != nil {
		return false, err
	}

	return total > 0, nil
}

// LoadRechargeRecord loads recharge record information.
func (d *SQLDB) LoadRechargeRecord(orderID string) (*types.RechargeRecord, error) {
	var info types.RechargeRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE order_id=?", rechargeRecordTable)
	err := d.db.Get(&info, query, orderID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// LoadRechargeRecords loads recharge records with a specific state.
func (d *SQLDB) LoadRechargeRecords(state types.RechargeState) ([]*types.RechargeRecord, error) {
	var infos []*types.RechargeRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE state=? ", rechargeRecordTable)

	err := d.db.Select(&infos, query, state)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

// SaveWithdrawInfoAndUserBalance saves withdraw information and updates user balance.
func (d *SQLDB) SaveWithdrawInfoAndUserBalance(rInfo *types.WithdrawRecord, balance, oldBalance string) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Errorf("SaveWithdrawInfoAndUserBalance Rollback err:%s", err.Error())
		}
	}()

	query := fmt.Sprintf(
		`INSERT INTO %s (order_id, value, state, withdraw_addr, withdraw_hash,  user_id) 
		        VALUES (:order_id,  :value, :state, :withdraw_addr, :withdraw_hash, :user_id)`, withdrawRecordTable)
	_, err = tx.NamedExec(query, rInfo)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`UPDATE %s SET balance=? WHERE user_id=? AND balance=?`, userTable)
	_, err = tx.Exec(query, balance, rInfo.UserID, oldBalance)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// LoadWithdrawRecord loads withdraw record information.
func (d *SQLDB) LoadWithdrawRecord(orderID string) (*types.WithdrawRecord, error) {
	var info types.WithdrawRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE order_id=?", withdrawRecordTable)
	err := d.db.Get(&info, query, orderID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// UpdateWithdrawRecord updates withdraw record information.
func (d *SQLDB) UpdateWithdrawRecord(info *types.WithdrawRecord, newState types.WithdrawState) error {
	query := fmt.Sprintf(`UPDATE %s SET state=?, value=?, done_time=NOW(), 
	     withdraw_hash=?, executor=? WHERE order_id=? AND state=?`, withdrawRecordTable)
	_, err := d.db.Exec(query, newState, info.Value, info.WithdrawHash, info.Executor, info.OrderID, info.State)

	return err
}

// LoadWithdrawRecords loads withdraw records with optional filters.
func (d *SQLDB) LoadWithdrawRecords(limit, page int64, statuses []types.WithdrawState, userID, startDate, endDate string) (*types.GetWithdrawResponse, error) {
	out := new(types.GetWithdrawResponse)

	whereStr := ""
	if userID != "" {
		whereStr = "AND user_id='" + userID + "'"
	}

	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			log.Errorf("Parse time err:%s", err.Error())
		} else {
			whereStr = whereStr + "AND created_time>='" + t.String() + "'"
		}
	}

	if endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			log.Errorf("Parse time err:%s", err.Error())
		} else {
			whereStr = whereStr + "AND created_time<='" + t.String() + "'"
		}
	}

	var infos []*types.WithdrawRecord
	lQuery := fmt.Sprintf("SELECT * FROM %s WHERE state in (?) %s order by created_time desc LIMIT ? OFFSET ?", withdrawRecordTable, whereStr)
	if limit > loadWithdrawRecordsDefaultLimit {
		limit = loadWithdrawRecordsDefaultLimit
	}
	lQuery, lArgs, err := sqlx.In(lQuery, statuses, limit, page*limit)
	if err != nil {
		return nil, err
	}
	lQuery = d.db.Rebind(lQuery)

	err = d.db.Select(&infos, lQuery, lArgs...)
	if err != nil {
		return nil, err
	}

	cQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE state in (?) %s", withdrawRecordTable, whereStr)
	var count int

	cQuery, cArgs, err := sqlx.In(cQuery, statuses)
	if err != nil {
		return nil, err
	}
	cQuery = d.db.Rebind(cQuery)

	err = d.db.Get(&count, cQuery, cArgs...)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

// LoadWithdrawRecordRows loads withdraw record rows with optional filters.
func (d *SQLDB) LoadWithdrawRecordRows(statuses []types.WithdrawState, userID, startDate, endDate string) (*sqlx.Rows, error) {
	whereStr := ""
	if userID != "" {
		whereStr = "AND user_id='" + userID + "'"
	}

	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			log.Errorf("Parse time err:%s", err.Error())
		} else {
			whereStr = whereStr + "AND created_time>='" + t.String() + "'"
		}
	}

	if endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			log.Errorf("Parse time err:%s", err.Error())
		} else {
			whereStr = whereStr + "AND created_time<='" + t.String() + "'"
		}
	}

	lQuery := fmt.Sprintf("SELECT * FROM %s WHERE state in (?) %s order by created_time desc ", withdrawRecordTable, whereStr)
	query, args, err := sqlx.In(lQuery, statuses)
	if err != nil {
		return nil, err
	}

	query = d.db.Rebind(query)
	return d.db.QueryxContext(context.Background(), query, args...)
}

// LoadWithdrawRecordsByUser loads withdraw records for a specific user with pagination.
func (d *SQLDB) LoadWithdrawRecordsByUser(userID string, limit, page int64) (*types.GetWithdrawResponse, error) {
	out := new(types.GetWithdrawResponse)

	var infos []*types.WithdrawRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? order by created_time desc LIMIT ? OFFSET ?", withdrawRecordTable)
	if limit > loadWithdrawRecordsDefaultLimit {
		limit = loadWithdrawRecordsDefaultLimit
	}

	err := d.db.Select(&infos, query, userID, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=?", withdrawRecordTable)
	var count int
	err = d.db.Get(&count, countQuery, userID)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

// LoadRechargeRecordsByUser loads recharge records for a specific user with pagination.
func (d *SQLDB) LoadRechargeRecordsByUser(userID string, limit, page int64) (*types.RechargeResponse, error) {
	out := new(types.RechargeResponse)

	var infos []*types.RechargeRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? order by created_time desc LIMIT ? OFFSET ?", rechargeRecordTable)
	if limit > loadRechargeRecordsDefaultLimit {
		limit = loadRechargeRecordsDefaultLimit
	}

	err := d.db.Select(&infos, query, userID, limit, page*limit)
	if err != nil {
		return nil, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=?", rechargeRecordTable)
	var count int
	err = d.db.Get(&count, countQuery, userID)
	if err != nil {
		return nil, err
	}

	out.Total = count
	out.List = infos

	return out, nil
}

// LoadWithdrawRecordsByUserAndState loads withdraw records for a specific user and state.
func (d *SQLDB) LoadWithdrawRecordsByUserAndState(userID string, state types.WithdrawState) ([]*types.WithdrawRecord, error) {
	var infos []*types.WithdrawRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND state=?", withdrawRecordTable)

	err := d.db.Select(&infos, query, userID, state)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
