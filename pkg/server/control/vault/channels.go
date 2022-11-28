package vault

import (
	"context"
	"database/sql"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	channelv3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/pkg/errors"
)

// CreateChannelStatus
func CreateChannelStatus(ctx context.Context, db *sql.DB, dialect excute.SqlExcutor, uuid, message string, created time.Time, max_count uint) error {
	var err error
	// created := time.Now()
	channel_status := channelv3.ChannelStatus{
		Uuid:    uuid,
		Created: created,
		Message: message,
	}

	channel_status_opt_cond := stmt.And(
		stmt.Equal("uuid", uuid),
	)

	var channel_opt channelv3.ChannelStatusOption
	// valid; exist channel
	err = dialect.QueryRows(channel_opt.TableName(), channel_opt.ColumnNames(), channel_status_opt_cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			err := channel_opt.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to query to channel opt")
	}

	if channel_opt.StatusMaxCount == 0 {
		channel_opt.StatusMaxCount = max_count // default channel status max_count
	}

	// rotation; find
	rotation_cond := stmt.And(
		stmt.Equal("uuid", uuid),
	)
	rotation_order := stmt.Desc("created")
	rotation_limit := stmt.Limit(int(channel_opt.StatusMaxCount)-1, 2)
	rotation_columns := []string{
		"uuid",
		"created",
	}

	uuids := make([]string, 0, state.ENV__INIT_SLICE_CAPACITY__())
	createds := make([]vanilla.NullTime, 0, state.ENV__INIT_SLICE_CAPACITY__())
	err = dialect.QueryRows(channel_status.TableName(), rotation_columns, rotation_cond, rotation_order, rotation_limit)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			var uuid string
			var created vanilla.NullTime

			err = scan.Scan(&uuid, &created)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			uuids = append(uuids, uuid)
			createds = append(createds, created)

			return err
		})
	if err != nil {
		return err
	}

	// insert
	// err = vanilla.Scope(db, func(tx *sql.Tx) (err error) {
	err = func() (err error) {
		affected, _, err := dialect.Insert(channel_status.TableName(), channel_status.ColumnNames(), channel_status.Values())(ctx, db)
		if err != nil {
			return
		}
		if affected == 0 {
			return errors.Errorf("no affected")
		}

		return
	}()
	if err != nil {
		return err
	}

	err = func() (err error) {
		// rotation; remove
		for i := range uuids {
			uuid, created := uuids[i], createds[i]
			rm_cond := stmt.And(
				stmt.Equal("uuid", uuid),
				stmt.Equal("created", created),
			)
			_, err = dialect.Delete(channel_status.TableName(), rm_cond)(ctx, db)
			if err != nil {
				return
			}
		}

		return
	}()
	if err != nil {
		return err
	}

	return nil
}

func GetManagedChannel(ctx context.Context, db *sql.DB, dialect excute.SqlExcutor, uuid string, tenant_hash string) (*channelv3.HttpRsp_ManagedChannel, error) {
	var err error
	var rst = new(channelv3.HttpRsp_ManagedChannel)

	channel_cond := stmt.And(
		stmt.Equal("uuid", uuid),
	)

	scoped_channel_cond := stmt.And(
		channel_cond,
		stmt.IsNull("deleted"),
	)

	// ManagedChannel
	var found bool
	table_channel := channelv3.TableNameWithTenant_ManagedChannel(tenant_hash)
	var channel channelv3.ManagedChannel
	err = dialect.QueryRows(table_channel, channel.ColumnNames(), scoped_channel_cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, i int) error {

			err := channel.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			found = true
			rst.ManagedChannel = channel

			return err
		})
	if err != nil {
		return rst, errors.Wrapf(err, "failed to get %v", channel.TableName())
	}
	if !found {
		return nil, nil
	}

	// ChannelStatusOption
	var status_option channelv3.ChannelStatusOption
	err = dialect.QueryRows(status_option.TableName(), status_option.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			err := status_option.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rst.StatusOption = status_option.ChannelStatusOption_property

			return err
		})
	if err != nil {
		return rst, errors.Wrapf(err, "failed to get %v", status_option.TableName())
	}

	// Format
	var format channelv3.Format
	err = dialect.QueryRows(format.TableName(), format.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			err := format.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rst.Format = format.Format_property

			return err
		})
	if err != nil {
		return rst, errors.Wrapf(err, "failed to get %v", format.TableName())
	}

	// NotifierEdge
	var set_edge = map[channelv3.NotifierEdge]struct{}{}
	var edge channelv3.NotifierEdge
	err = dialect.QueryRows(edge.TableName(), edge.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
		func(scan excute.Scanner, _ int) error {
			err := edge.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			set_edge[edge] = struct{}{}

			return err
		})
	if err != nil {
		return rst, errors.Wrapf(err, "failed to get %v", edge.TableName())
	}

	for edge := range set_edge {
		edge_options, err := GetChannelNotifierEdge(ctx, db, dialect, edge)
		if err != nil {
			return rst, errors.Wrapf(err, "failed to get channel notifier edge")
		}

		rst.Notifiers.NotifierEdge_property = edge_options.NotifierEdge_property
		rst.Notifiers.Console = edge_options.Console
		rst.Notifiers.Webhook = edge_options.Webhook
		rst.Notifiers.RabbitMq = edge_options.RabbitMq
		rst.Notifiers.Slackhook = edge_options.Slackhook
	}

	return rst, nil
}

func GetChannelNotifierEdge(ctx context.Context, db *sql.DB, dialect excute.SqlExcutor, edge channelv3.NotifierEdge) (*channelv3.HttpRsp_ManagedChannel_NotifierEdge, error) {
	var err error
	var rsp = new(channelv3.HttpRsp_ManagedChannel_NotifierEdge)

	rsp.NotifierEdge = edge

	channel_cond := stmt.And(
		stmt.Equal("uuid", edge.Uuid),
	)

	if edge.NotifierType == channelv3.NotifierTypeConsole {
		// NotifierConsole
		var notifier channelv3.NotifierConsole
		err = dialect.QueryRows(notifier.TableName(), notifier.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
			func(scan excute.Scanner, _ int) error {
				err := notifier.Scan(scan)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}

				rsp.Console = &notifier.NotifierConsole_property

				return err
			})

		if err != nil {
			return rsp, errors.Wrapf(err, "failed to get a notifier(console)")
		}

	}

	if edge.NotifierType == channelv3.NotifierTypeWebhook {
		// NotifierWebhook
		var notifier channelv3.NotifierWebhook
		err = dialect.QueryRows(notifier.TableName(), notifier.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
			func(scan excute.Scanner, _ int) error {
				err := notifier.Scan(scan)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}

				rsp.Webhook = &notifier.NotifierWebhook_property

				return err
			})
		if err != nil {
			return rsp, errors.Wrapf(err, "failed to get a notifier(webhook)")
		}

	}

	if edge.NotifierType == channelv3.NotifierTypeRabbitmq {
		// NotifierRabbitMq
		var notifier channelv3.NotifierRabbitMq
		err = dialect.QueryRows(notifier.TableName(), notifier.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
			func(scan excute.Scanner, _ int) error {
				err := notifier.Scan(scan)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}

				rsp.RabbitMq = &notifier.NotifierRabbitMq_property

				return err
			})
		if err != nil {
			return rsp, errors.Wrapf(err, "failed to get a notifier(rabbitmq)")
		}

	}

	if edge.NotifierType == channelv3.NotifierTypeSlackhook {
		// NotifierSlackhook
		var notifier channelv3.NotifierSlackhook
		err = dialect.QueryRows(notifier.TableName(), notifier.ColumnNames(), channel_cond, nil, nil)(ctx, db)(
			func(scan excute.Scanner, _ int) error {
				err := notifier.Scan(scan)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}

				rsp.Slackhook = &notifier.NotifierSlackhook_property

				return err
			})

		if err != nil {
			return rsp, errors.Wrapf(err, "failed to get a notifier(slackhook)")
		}

	}

	return rsp, nil
}

func CheckManagedChannel(ctx context.Context, db excute.Preparer, dialect excute.SqlExcutor, uuid string, tenant_hash string) error {
	// check channel
	var channel channelv3.ManagedChannel
	channel.Uuid = uuid
	channel_cond := stmt.And(
		stmt.Equal("uuid", channel.Uuid),
		stmt.IsNull("deleted"),
	)
	table_channel := channelv3.TableNameWithTenant_ManagedChannel(tenant_hash)

	exist, err := dialect.Exist(table_channel, channel_cond)(ctx, db)
	if err != nil {
		return errors.Wrapf(err, "failed to check a channel")
	}
	if !exist {
		return errors.WithStack(database.ErrorRecordWasNotFound)
	}

	return nil
}
