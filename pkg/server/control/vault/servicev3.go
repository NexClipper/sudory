package vault

import (
	"context"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	service "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	"github.com/pkg/errors"
)

func GetService(
	ctx context.Context,
	db excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, uuid string,
) (*service.Service, error) {
	var table service.Service

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, stmt.Equal("uuid", uuid))

	cond := stmt.And(args...)

	var record *service.Service
	err := dialect.QueryRows(table.TableName(), table.ColumnNames(), cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			var service service.Service
			err := service.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			if record == nil {
				record = &service
			}

			// find latest record by created
			if record.Timestamp.Before(service.Timestamp) {
				record = &service
			}

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service")
		return record, err
	}
	if record == nil {
		err = errors.New("could not found service record")
		return record, err
	}

	return record, err
}

func GetServicesPolling(
	ctx context.Context,
	db excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, polling_offset vanilla.NullTime,
) ([]service.Service_polling, error) {

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	if polling_offset.Valid {
		args = append(args, stmt.GTE("created", polling_offset))
	}

	if len(args) == 0 {
		return nil, errors.New("need more conditon")
	}

	cond := stmt.And(args...)

	// limit := vanilla.Limit(math.MaxInt8)

	var recordSet_ = make(map[string]service.Service_polling)
	var record service.Service_polling
	err := dialect.QueryRows(record.TableName(), record.ColumnNames(), cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, i int) error {
			err := record.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			if _, ok := recordSet_[record.Uuid]; !ok {
				recordSet_[record.Uuid] = record
			}

			// find latest record by created
			if recordSet_[record.Uuid].Timestamp.Before(record.Timestamp) {
				recordSet_[record.Uuid] = record
			}

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service")
		return nil, err
	}

	var records = make([]service.Service_polling, 0, len(recordSet_))
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

	return records, nil
}

func GetServiceStep(
	ctx context.Context,
	db excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, uuid string, sequence int,
) (*service.ServiceStep, error) {

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, stmt.Equal("uuid", uuid))
	args = append(args, stmt.Equal("seq", sequence))

	cond := stmt.And(args...)

	var record *service.ServiceStep
	var table service.ServiceStep
	err := dialect.QueryRows(table.TableName(), table.ColumnNames(), cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			var step service.ServiceStep
			err := step.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			if record == nil {
				record = &step
			}

			// find latest record by created
			if record.Timestamp.Before(step.Timestamp) {
				record = &step
			}

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return record, err
	}
	if record == nil {
		err = errors.New("could not found service step record")
		return record, err
	}

	return record, nil
}

func GetServiceSteps(
	ctx context.Context,
	db excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, uuid string,
) ([]service.ServiceStep, error) {

	args := make([]map[string]interface{}, 0, 2)
	if 0 < len(cluster_uuid) {
		args = append(args, stmt.Equal("cluster_uuid", cluster_uuid))
	}
	args = append(args, stmt.Equal("uuid", uuid))

	cond := stmt.And(args...)

	var recordSet = make(map[int]service.ServiceStep)
	var record service.ServiceStep
	err := dialect.QueryRows(record.TableName(), record.ColumnNames(), cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			err := record.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			if _, ok := recordSet[record.Sequence]; !ok {
				recordSet[record.Sequence] = record
			}

			// find latest record by created
			if recordSet[record.Sequence].Timestamp.Before(record.Timestamp) {
				recordSet[record.Sequence] = record
			}

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return nil, err
	}

	var records = make([]service.ServiceStep, 0, len(recordSet))
	for key := range recordSet {
		records = append(records, recordSet[key])
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Sequence < records[j].Sequence
	})

	return records, nil
}

func GetServiceResult(
	ctx context.Context,
	db excute.Preparer, dialect excute.SqlExcutor,
	cond stmt.ConditionStmt, order stmt.OrderStmt, page stmt.PaginationStmt,
) (*service.ServiceResult, error) {

	var record *service.ServiceResult
	var table service.ServiceResult
	err := dialect.QueryRows(table.TableName(), table.ColumnNames(), cond, order, page)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			var result service.ServiceResult
			err := result.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			if record == nil {
				record = &result
			}

			// find latest record by created
			if record.Timestamp.Before(result.Timestamp) {
				record = &result
			}

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service result")
		return record, err
	}
	if record == nil {
		err = errors.New("could not found service result record")
		return record, err
	}

	return record, nil
}

type Table interface {
	TableName() string
	Values() []interface{}
	ColumnNames() []string
}

func SaveMultiTable(tx excute.Preparer, dialect excute.SqlExcutor, tables []Table) error {
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

	affected, _, err := dialect.Insert(tablename, columnnames, values...)(context.Background(), tx)
	if err != nil {
		return errors.Wrapf(err, "could not save")
	}

	if affected == 0 {
		return errors.New("no affected")
	}

	return nil
}
