package vault

import (
	"context"
	"database/sql"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	v3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	"github.com/pkg/errors"
)

type servicev3 struct{}

var (
	Servicev3 = new(servicev3)
)

func (servicev3) GetService(
	ctx context.Context,
	db vanilla.Preparer,
	cluster_uuid string, uuid string,
) (record *v3.Service, err error) {
	var table v3.Service

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, vanilla.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, vanilla.Equal("uuid", uuid))

	cond := vanilla.And(args[0], args[1:]...)

	err = vanilla.Stmt.Select(table.TableName(), table.ColumnNames(), cond.Parse(), nil, nil).
		QueryRowsContext(ctx, db)(func(scan vanilla.Scanner, _ int) error {

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
	db vanilla.Preparer,
	cluster_uuid string, polling_offset vanilla.NullTime,
) (records []v3.Service_polling, err error) {
	var record v3.Service_polling
	recordSet_ := make(map[string]v3.Service_polling)

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, vanilla.Equal("cluster_uuid", cluster_uuid))
	}
	if polling_offset.Valid {
		args = append(args, vanilla.GreaterThanEqual("created", polling_offset))
	}

	if len(args) == 0 {
		err = errors.New("need more conditon")
		return
	}

	cond := vanilla.And(args[0], args[1:]...)

	// limit := vanilla.Limit(math.MaxInt8)

	err = vanilla.Stmt.Select(record.TableName(), record.ColumnNames(), cond.Parse(), nil, nil).
		QueryRowsContext(ctx, db)(func(scan vanilla.Scanner, i int) error {

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
		if records[i].Priority > records[j].Priority {
			return true
		} else if records[i].Priority < records[j].Priority {
			return false
		} else {
			return records[i].Created.Before(records[j].Created)
		}
	})

	return
}

func (servicev3) GetServiceStep(
	ctx context.Context,
	db vanilla.Preparer,
	cluster_uuid string, uuid string, sequence int,
) (record *v3.ServiceStep, err error) {
	var table v3.ServiceStep

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, vanilla.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, vanilla.Equal("uuid", uuid))
	args = append(args, vanilla.Equal("seq", sequence))

	cond := vanilla.And(args[0], args[1:]...)

	err = vanilla.Stmt.Select(table.TableName(), table.ColumnNames(), cond.Parse(), nil, nil).
		QueryRowsContext(ctx, db)(func(scan vanilla.Scanner, _ int) error {

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
	db vanilla.Preparer,
	cluster_uuid string, uuid string,
) (records map[string][]v3.ServiceStep, err error) {

	recordSet := make(map[string]map[int]v3.ServiceStep)
	var record v3.ServiceStep

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, vanilla.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, vanilla.Equal("uuid", uuid))

	cond := vanilla.And(args[0], args[1:]...)

	err = vanilla.Stmt.Select(record.TableName(), record.ColumnNames(), cond.Parse(), nil, nil).
		QueryRowsContext(ctx, db)(func(scan vanilla.Scanner, _ int) error {

		err = record.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service step")
		}
		// init sub record set
		if recordSet[record.Uuid] == nil {
			recordSet[record.Uuid] = make(map[int]v3.ServiceStep)
		}

		if _, ok := recordSet[record.Uuid][record.Sequence]; !ok {
			recordSet[record.Uuid][record.Sequence] = record
		}

		// find latest record by created
		if recordSet[record.Uuid][record.Sequence].Timestamp.Before(record.Timestamp) {
			recordSet[record.Uuid][record.Sequence] = record
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return
	}

	records = make(map[string][]v3.ServiceStep)
	for uuid, stepSet := range recordSet {
		if records[uuid] == nil {
			records[uuid] = make([]v3.ServiceStep, 0, len(stepSet))
		}
		for _, setp := range stepSet {
			records[uuid] = append(records[uuid], setp)
		}
	}

	for uuid := range records {
		sort.Slice(records[uuid], func(i, j int) bool {
			return records[uuid][i].Sequence < records[uuid][j].Sequence
		})
	}

	return
}

func (servicev3) GetServiceResult(
	ctx context.Context,
	db vanilla.Preparer,
	cond *prepare.Condition, order *prepare.Orders, page *prepare.Pagination,
) (record *v3.ServiceResult, err error) {
	var table v3.ServiceResult

	err = vanilla.Stmt.Select(table.TableName(), table.ColumnNames(), cond, order, page).
		QueryRowsContext(ctx, db)(func(scan vanilla.Scanner, _ int) error {

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

func SaveMultiTable(tx *sql.Tx, tables []Table) error {
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

	builder, err := vanilla.Stmt.Insert(tablename, columnnames, values...)
	if err != nil {
		return errors.Wrapf(err, "could not build a insert statement")
	}

	affected, err := builder.Exec(tx)
	if err != nil {
		return errors.Wrapf(err, "exec insert statement")
	}
	if affected == 0 {
		return errors.New("no affected")
	}

	return nil
}
