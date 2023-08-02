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

	return s, nil
}

const (
	// Database table names.
	assetRecordTable = "asset_record"

	// Default limits for loading table entries.
	loadNodeInfosDefaultLimit           = 100
	loadValidationResultsDefaultLimit   = 100
	loadAssetRecordsDefaultLimit        = 100
	loadExpiredAssetRecordsDefaultLimit = 100
	loadWorkloadDefaultLimit            = 100
	loadReplicaEventDefaultLimit        = 100
	loadRetrieveDefaultLimit            = 100
	loadReplicaDefaultLimit             = 100
)

// InitTables initializes data tables.
func InitTables(d *SQLDB) error {
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
	tx.MustExec(fmt.Sprintf(cAssetRecordTable, assetRecordTable))

	return tx.Commit()
}
