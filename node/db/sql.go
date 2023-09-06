package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("db")

// SQLDB represents a scheduler sql database.
type SQLDB struct {
	db *sqlx.DB
}

// NewSQLDB creates a new database connection using the given MySQL connection string.
// The function returns a SQLDB pointer or an error if the connection failed.
func NewSQLDB(path string) (*SQLDB, error) {
	path = fmt.Sprintf("%s?parseTime=true&loc=Local", path)

	client, err := sqlx.Open("mysql", path)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(); err != nil {
		return nil, err
	}

	s := &SQLDB{client}
	s.initTables()

	return s, nil
}

const (
	// Database table names.
	orderRecordTable      = "order_record"
	rechargeRecordTable   = "recharge_record"
	withdrawRecordTable   = "withdraw_record"
	userInstancesTable    = "user_instances_details"
	configTable           = "config"
	userTable             = "user_info"
	adminTable            = "admin_info"
	rechargeAddressTable  = "recharge_address"
	instanceBaseInfoTable = "instance_base_info"
	instanceRefundTable   = "instance_refund"
	invitationTable       = "invitation"
	providerInfoTable     = "provider_info"
	accessKeyTable        = "provider_access_info"

	// Default limits for loading table entries.
	loadOrderRecordsDefaultLimit    = 1000
	loadRechargeRecordsDefaultLimit = 1000
	loadWithdrawRecordsDefaultLimit = 1000
	loadAddressesDefaultLimit       = 1000
	loadInstancesDefaultLimit       = 100
)

// initTables initializes data tables.
func (d *SQLDB) initTables() error {
	// init table
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Errorf("InitTables Rollback err:%s", err.Error())
		}
	}()

	// Execute table creation statements
	tx.MustExec(fmt.Sprintf(cOrderRecordTable, orderRecordTable))
	tx.MustExec(fmt.Sprintf(cInstanceDetailsTable, userInstancesTable))
	tx.MustExec(fmt.Sprintf(cRechargeTable, rechargeRecordTable))
	tx.MustExec(fmt.Sprintf(cWithdrawTable, withdrawRecordTable))
	tx.MustExec(fmt.Sprintf(cConfigTable, configTable))
	tx.MustExec(fmt.Sprintf(cUserTable, userTable))
	tx.MustExec(fmt.Sprintf(cRechargeAddressTable, rechargeAddressTable))
	tx.MustExec(fmt.Sprintf(cAdminTable, adminTable))
	tx.MustExec(fmt.Sprintf(cInstanceDefaultTable, instanceBaseInfoTable))
	tx.MustExec(fmt.Sprintf(cInstanceRefundTable, instanceRefundTable))
	tx.MustExec(fmt.Sprintf(cInvitationTable, invitationTable))
	tx.MustExec(fmt.Sprintf(cProviderInfoTable, providerInfoTable))
	tx.MustExec(fmt.Sprintf(cAccessKeyTable, accessKeyTable))

	return tx.Commit()
}
