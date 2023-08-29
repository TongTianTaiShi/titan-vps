package orders

import (
	"bytes"
	"context"
	"strings"

	"github.com/LMF709268224/titan-vps/api/types"
	"github.com/LMF709268224/titan-vps/node/db"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/jmoiron/sqlx"
)

// Datastore represents the order datastore
type Datastore struct {
	orderDB *db.SQLDB
}

// NewDatastore creates a new datastore
func NewDatastore(db *db.SQLDB) *Datastore {
	return &Datastore{
		orderDB: db,
	}
}

// Close closes the order datastore
func (d *Datastore) Close() error {
	return nil
}

func trimPrefix(key datastore.Key) string {
	return strings.Trim(key.String(), "/")
}

// Get retrieves data from the datastore
func (d *Datastore) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	cInfo, err := d.orderDB.LoadOrderRecord(trimPrefix(key), orderTimeoutMinute)
	if err != nil {
		return nil, err
	}

	order := orderInfoFrom(cInfo)

	valueBuf := new(bytes.Buffer)
	if err := order.MarshalCBOR(valueBuf); err != nil {
		return nil, err
	}

	return valueBuf.Bytes(), nil
}

// Has  checks if the key exists in the datastore
func (d *Datastore) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	return d.orderDB.OrderExists(trimPrefix(key))
}

// GetSize gets the data size from the datastore
func (d *Datastore) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	return d.orderDB.LoadOrderCount()
}

// Query queries order records from the datastore
func (d *Datastore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	var rows *sqlx.Rows
	var err error

	rows, err = d.orderDB.LoadAllOrderRecords(ActiveStates, orderTimeoutMinute)
	if err != nil {
		log.Errorf("LoadAllOrderRecords :%s", err.Error())
		return nil, err
	}
	defer rows.Close()

	re := make([]query.Entry, 0)
	// loading orders to local
	for rows.Next() {
		cInfo := &types.OrderRecord{}
		err = rows.StructScan(cInfo)
		if err != nil {
			log.Errorf("StructScan err: %s", err.Error())
			continue
		}

		order := orderInfoFrom(cInfo)
		valueBuf := new(bytes.Buffer)
		if err = order.MarshalCBOR(valueBuf); err != nil {
			log.Errorf("order marshal cbor: %s", err.Error())
			continue
		}

		prefix := "/"
		entry := query.Entry{
			Key: prefix + order.OrderID.String(), Size: len(valueBuf.Bytes()),
		}

		if !q.KeysOnly {
			entry.Value = valueBuf.Bytes()
		}

		re = append(re, entry)
	}

	r := query.ResultsWithEntries(q, re)
	r = query.NaiveQueryApply(q, r)

	return r, nil
}

// Put update order record info
func (d *Datastore) Put(ctx context.Context, key datastore.Key, value []byte) error {
	aInfo := &OrderInfo{}
	if err := aInfo.UnmarshalCBOR(bytes.NewReader(value)); err != nil {
		return err
	}

	aInfo.OrderID = OrderHash(trimPrefix(key))

	return d.orderDB.SaveOrderInfo(aInfo.ToOrderRecord())
}

// Delete delete order record info (This func has no place to call it)
func (d *Datastore) Delete(ctx context.Context, key datastore.Key) error {
	return nil
}

// Sync sync
func (d *Datastore) Sync(ctx context.Context, prefix datastore.Key) error {
	return nil
}

// Batch batch
func (d *Datastore) Batch(ctx context.Context) (datastore.Batch, error) {
	return nil, nil
}

var _ datastore.Batching = (*Datastore)(nil)
