package vault

import (
	"context"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	service "github.com/NexClipper/sudory/pkg/server/model/service/v4"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/pkg/errors"
)

type ServicePolling struct {
	Service    service.Service
	LastStatus service.ServiceStatus
}

func GetServicesPolling_v4(
	ctx context.Context,
	db excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, polling_offset vanilla.NullTime,
) ([]ServicePolling, error) {

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

	var records = make([]ServicePolling, 0, state.ENV__INIT_SLICE_CAPACITY__())
	var recordSet = make(map[string]int)
	var serv service.Service
	err := dialect.QueryRows(serv.TableName(), serv.ColumnNames(), cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, i int) error {

			var serv service.Service
			err := serv.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			defaultStatus := service.ServiceStatus{}
			defaultStatus.PartitionDate = serv.PartitionDate
			defaultStatus.Created = serv.Created
			defaultStatus.ClusterUuid = serv.ClusterUuid
			defaultStatus.Uuid = serv.Uuid
			defaultStatus.StepMax = serv.StepMax

			records = append(records, ServicePolling{Service: serv, LastStatus: defaultStatus})
			recordSet[serv.Uuid] = i

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service")
		return nil, err
	}

	var status service.ServiceStatus
	err = dialect.QueryRows(status.TableName(), status.ColumnNames(), cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {

			var status service.ServiceStatus
			err := status.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			i, ok := recordSet[status.Uuid]
			if ok && records[i].LastStatus.Created.Before(status.Created) {
				records[i].LastStatus = status
			}

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service status")
		return nil, err
	}

	// sort by priority, created
	sort.Slice(records, func(i, j int) bool {
		if records[i].Service.Priority != records[j].Service.Priority {
			return records[i].Service.Priority > records[j].Service.Priority
		}

		return records[i].Service.Created.Before(records[j].Service.Created)
	})

	return records, nil
}
