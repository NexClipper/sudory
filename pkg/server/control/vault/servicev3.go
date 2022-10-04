package vault

import (
	"context"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	v3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	"github.com/pkg/errors"
)

type servicev3 struct{}

var (
	Servicev3 = new(servicev3)
)

func (servicev3) GetService(
	ctx context.Context,
	db stmtex.Preparer, dialect string,
	cluster_uuid string, uuid string,
) (record *v3.Service, err error) {
	var table v3.Service

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, stmt.Equal("uuid", uuid))

	cond := stmt.And(args...)

	err = stmtex.Select(table.TableName(), table.ColumnNames(), cond, nil, nil).
		QueryRowsContext(ctx, db, dialect)(func(scan stmtex.Scanner, _ int) error {

		tmp := new(v3.Service)
		err = tmp.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service")
		}

		if record == nil {
			record = tmp
		}

		// find latest record by created
		if record.Timestamp.Before(tmp.Timestamp) {
			record = tmp
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service")
		return
	}
	if record == nil {
		err = errors.New("could not found service record")
		return
	}

	return
}

func (servicev3) GetServicesPolling(
	ctx context.Context,
	db stmtex.Preparer, dialect string,
	cluster_uuid string, polling_offset vanilla.NullTime,
) (records []v3.Service_polling, err error) {
	var record v3.Service_polling
	recordSet_ := make(map[string]v3.Service_polling)

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	if polling_offset.Valid {
		args = append(args, stmt.GTE("created", polling_offset))
	}

	if len(args) == 0 {
		err = errors.New("need more conditon")
		return
	}

	cond := stmt.And(args...)

	// limit := vanilla.Limit(math.MaxInt8)

	err = stmtex.Select(record.TableName(), record.ColumnNames(), cond, nil, nil).
		QueryRowsContext(ctx, db, dialect)(func(scan stmtex.Scanner, i int) error {

		err = record.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service")
		}

		if _, ok := recordSet_[record.Uuid]; !ok {
			recordSet_[record.Uuid] = record
		}

		// find latest record by created
		if recordSet_[record.Uuid].Timestamp.Before(record.Timestamp) {
			recordSet_[record.Uuid] = record
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service")
		return
	}

	records = make([]v3.Service_polling, 0, len(recordSet_))
	for _, service := range recordSet_ {
		records = append(records, service)
	}
	// sort by priority, created
	sort.Slice(records, func(i, j int) bool {
		if records[i].Priority != records[j].Priority {
			return records[i].Priority > records[j].Priority
		}

		return records[i].Created.Before(records[j].Created)
	})

	return
}

func (servicev3) GetServiceStep(
	ctx context.Context,
	db stmtex.Preparer, dialect string,
	cluster_uuid string, uuid string, sequence int,
) (record *v3.ServiceStep, err error) {
	var table v3.ServiceStep

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, stmt.Equal("uuid", uuid))
	args = append(args, stmt.Equal("seq", sequence))

	cond := stmt.And(args...)

	err = stmtex.Select(table.TableName(), table.ColumnNames(), cond, nil, nil).
		QueryRowsContext(ctx, db, dialect)(func(scan stmtex.Scanner, _ int) error {

		tmp := new(v3.ServiceStep)
		err = tmp.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service step")
		}

		if record == nil {
			record = tmp
		}

		// find latest record by created
		if record.Timestamp.Before(tmp.Timestamp) {
			record = tmp
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return
	}
	if record == nil {
		err = errors.New("could not found service step record")
		return
	}

	return
}

func (servicev3) GetServiceSteps(
	ctx context.Context,
	db stmtex.Preparer, dialect string,
	cluster_uuid string, uuid string,
) (records []v3.ServiceStep, err error) {

	recordSet := make(map[int]v3.ServiceStep)
	var record v3.ServiceStep

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, stmt.Equal("uuid", uuid))

	cond := stmt.And(args...)

	err = stmtex.Select(record.TableName(), record.ColumnNames(), cond, nil, nil).
		QueryRowsContext(ctx, db, dialect)(func(scan stmtex.Scanner, _ int) error {

		err = record.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service step")
		}

		if _, ok := recordSet[record.Sequence]; !ok {
			recordSet[record.Sequence] = record
		}

		// find latest record by created
		if recordSet[record.Sequence].Timestamp.Before(record.Timestamp) {
			recordSet[record.Sequence] = record
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return
	}

	records = make([]v3.ServiceStep, 0, len(recordSet))
	for key := range recordSet {
		records = append(records, recordSet[key])
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Sequence < records[j].Sequence
	})

	return
}

func (servicev3) GetServiceResult(
	ctx context.Context,
	db stmtex.Preparer, dialect string,
	cond stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt,
) (record *v3.ServiceResult, err error) {
	var table v3.ServiceResult

	err = stmtex.Select(table.TableName(), table.ColumnNames(), cond, order, page).
		QueryRowsContext(ctx, db, dialect)(func(scan stmtex.Scanner, _ int) error {

		tmp := new(v3.ServiceResult)
		err = tmp.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service result")
		}

		if record == nil {
			record = tmp
		}

		// find latest record by created
		if record.Timestamp.Before(tmp.Timestamp) {
			record = tmp
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service result")
		return
	}
	if record == nil {
		err = errors.New("could not found service result record")
		return
	}

	return
}

type Table interface {
	TableName() string
	Values() []interface{}
	ColumnNames() []string
}

func SaveMultiTable(tx stmtex.Preparer, dialect string, tables []Table) error {
	BuildInsertValues := func(tables []Table) [][]interface{} {
		values := make([][]interface{}, len(tables))
		for i, table := range tables {
			values[i] = make([]interface{}, 0, len(table.ColumnNames()))
			values[i] = append(values[i], table.Values()...)
		}
		return values
	}

	var (
		tablename   string
		columnnames []string
		values      [][]interface{}
	)

	if 0 < len(tables) {
		tablename = tables[0].TableName()
		columnnames = tables[0].ColumnNames()
		values = BuildInsertValues(tables)
	}

	if len(tablename) == 0 || len(columnnames) == 0 || len(values) == 0 {
		// nothing to do
		return nil
	}

	affected, _, err := stmtex.Insert(tablename, columnnames, values...).
		Exec(tx, dialect)
	if err != nil {
		return errors.Wrapf(err, "could not save")
	}

	if affected == 0 {
		return errors.New("no affected")
	}

	return nil
}
