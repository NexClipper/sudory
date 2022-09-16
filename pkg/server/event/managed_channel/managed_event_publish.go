package managed_channel

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NexClipper/logger"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/pkg/errors"
)

var InvokeByChannelUuid func(tenant_hash string, channel_uuid string, v map[string]interface{}) = func(tenant_hash string, channel_uuid string, v map[string]interface{}) {}

var InvokeByEventCategory func(tenant_hash string, ec channelv3.EventCategory, v map[string]interface{}) = func(tenant_hash string, ec channelv3.EventCategory, v map[string]interface{}) {}

// var _ EventPublisher = (*ManagedChannel)(nil)

type Event struct {
	*sql.DB
	dialect string
	// *vanilla.SqlDbEx

	// HashsetEventNotifierMuxer
	EventNotifierMuxer

	ErrorHandlers         event.HashsetErrorHandlers
	NofitierErrorHandlers HashsetNofitierErrorHandler
}

func NewEvent(db *sql.DB, dialect string) *Event {

	me := Event{}
	me.DB = db
	me.dialect = dialect
	// me.SqlDbEx = &vanilla.SqlDbEx{DB: db}
	// me.HashsetEventNotifierMuxer = HashsetEventNotifierMuxer{}

	me.ErrorHandlers = event.HashsetErrorHandlers{}
	me.NofitierErrorHandlers = HashsetNofitierErrorHandler{}

	return &me
}

func (pub *Event) Dialect() string {
	return pub.dialect
}

func (pub *Event) SetEventNotifierMuxer(mux EventNotifierMuxer) {
	pub.EventNotifierMuxer = mux
}

// func (pub Event) EventNotifierMuxer() HashsetEventNotifierMuxer {
// 	return pub.HashsetEventNotifierMuxer
// }

func (pub Event) InvokeByChannelUuid(tenant_hash string, channel_uuid string, v map[string]interface{}) {
	clone := NewEvent(pub.DB, pub.Dialect())
	clone.ErrorHandlers = pub.ErrorHandlers
	clone.NofitierErrorHandlers = pub.NofitierErrorHandlers

	//build channel muxer
	if err := clone.BuildMuxerByChannelUuid(tenant_hash, channel_uuid); err != nil {
		clone.OnError(errors.Wrapf(err, "build muxer by channel_uuid"))
		return
	}

	for channel_uuid := range clone.EventNotifierMuxer.Notifiers() {
		//build formatter
		if err := clone.BuildChannelFormatter(channel_uuid); err != nil {
			clone.OnError(errors.Wrapf(err, "build managed event"))
			return
		}
	}

	//update message
	clone.Update(v)
}

func (pub Event) BuildChannelFormatter(channel_uuid string) (err error) {
	eq_uuid := stmt.Equal("uuid", channel_uuid)
	format := channelv3.Format{}
	err = stmtex.Select(format.TableName(), format.ColumnNames(), eq_uuid, nil, nil).
		QueryRows(pub, pub.Dialect())(func(scan stmtex.Scanner, _ int) (err error) {
		err = format.Scan(scan)
		if err != nil {
			return
		}
		switch format.FormatType {
		case channelv3.FormatTypeFields:
			formatter := &Formatter_fields{
				FormatData: format.FormatData,
			}
			pub.EventNotifierMuxer.Formatters().Add(channel_uuid, formatter)
		case channelv3.FormatTypeJq:
			formatter := &Formatter_jq{
				FormatData: format.FormatData,
			}
			pub.EventNotifierMuxer.Formatters().Add(channel_uuid, formatter)
		default:
			// do nothing
		}

		return
	})
	err = errors.Wrapf(err, "build channel formatter")

	return
}

func (pub Event) InvokeByEventCategory(tenant_hash string, ec channelv3.EventCategory, v map[string]interface{}) {
	clone := NewEvent(pub.DB, pub.Dialect())
	clone.ErrorHandlers = pub.ErrorHandlers
	clone.NofitierErrorHandlers = pub.NofitierErrorHandlers

	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(clone)

	//build channel muxer
	if err := clone.BuildMuxerByEventCategory(tenant_hash, ec); err != nil {
		clone.OnError(errors.Wrapf(err, "build muxer by channel_uuid"))
		return
	}

	for channel_uuid := range clone.EventNotifierMuxer.Notifiers() {
		//build formatter
		if err := clone.BuildChannelFormatter(channel_uuid); err != nil {
			clone.OnError(errors.Wrapf(err, "build managed event"))
			return
		}
	}
	//update message
	clone.Update(v)
}

var (
	DefaultErrorHandler = func(err error) {
		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = logs.KVL(
					"stack", s,
				)
			})
		})

		logger.Error(fmt.Errorf("%w%s", err, stack))
	}

	DefaultErrorHandler_notifier = func(pub *Event) func(notifier Notifier, err error) {
		return func(notifier Notifier, err error) {
			defer func() {
				r := recover()

				if r == nil {
					return
				}

				if err, ok := r.(error); ok {
					pub.OnError(errors.Wrapf(err, "recover notifier error handler"))
				} else {
					pub.OnError(errors.Errorf("notifier error handler recover='%+v'", r))
				}
			}()

			var stack string
			logs.CauseIter(err, func(err error) {
				logs.StackIter(err, func(s string) {
					stack = logs.KVL(
						"stack", s,
					)
				})
			})

			//이벤트 알림 상태 테이블에 에러 메시지 저장
			uuid := notifier.Uuid()
			created := time.Now()
			message := fmt.Sprintf("%s%s", err.Error(), stack)

			if err_ := vault.CreateChannelStatus(pub.DB, pub.Dialect(), uuid, message, created, globvar.Event.NofitierStatusRotateLimit()); err_ != nil {
				err_ = errors.Wrapf(err_, "failed to logging to channel status")
				pub.ErrorHandlers.OnError(err_)
			}
		}
	}
)

func (pub *Event) Close() {

}

func (pub *Event) OnError(err error) {
	pub.ErrorHandlers.OnError(err)
}
func (pub *Event) OnNotifierError(notifier Notifier, err error) {
	pub.NofitierErrorHandlers.OnError(notifier, err)
}

func (pub *Event) BuildMuxerByEventCategory(tenant_hash string, event_category channelv3.EventCategory) (err error) {
	var ctx = context.Background()
	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(pub)

	channel_cond := stmt.And(
		stmt.Equal("event_category", int(event_category)),
		stmt.IsNull("deleted"),
	)
	column_names := []string{"uuid"}

	// var channel channelv3.ManagedChannel

	channel_table := channelv3.TableNameWithTenant_ManagedChannel(tenant_hash)
	err = stmtex.Select(channel_table, column_names, channel_cond, nil, nil).
		QueryRowsContext(ctx, pub, pub.Dialect())(func(scan stmtex.Scanner, _ int) (err error) {

		var channel_uuid string
		err = scan.Scan(&channel_uuid)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}

		// get notifier edge with option
		var edge channelv3.NotifierEdge
		edge.Uuid = channel_uuid
		edge_cond := stmt.Equal("uuid", channel_uuid)

		err = stmtex.Select(edge.TableName(), edge.ColumnNames(), edge_cond, nil, nil).
			QueryRowsContext(ctx, pub, pub.Dialect())(
			func(scan stmtex.Scanner, _ int) error {
				err := edge.Scan(scan)
				if err != nil {
					return errors.Wrapf(err, "failed to scan")
				}

				edge_opt, err := vault.GetChannelNotifierEdge(ctx, pub.DB, pub.Dialect(), edge)
				if err != nil {
					return errors.Wrapf(err, "failed to get a NotifierEdge_option")
				}

				// valied notifier
				err = ValidNotifier(edge_opt)
				if err != nil {
					return errors.Wrapf(err, "valied notifier")
				}

				// notifier factory
				notifier, err := NotifierFactory(channel_uuid, edge_opt)
				if err != nil {
					return errors.Wrapf(err, "channel notifier factory")
				}

				// append notifier
				muxer.Notifiers().Add(channel_uuid, notifier)

				return nil
			})
		return
	})
	if err != nil {
		return
	}

	return nil
}

func (pub *Event) BuildMuxerByChannelUuid(tenant_hash string, channel_uuid string) error {
	var ctx = context.Background()
	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(pub)

	// get notifier edge with option

	var edge channelv3.NotifierEdge
	edge.Uuid = channel_uuid
	edge_cond := stmt.Equal("uuid", channel_uuid)
	edge_table := channelv3.TableNameWithTenant_NotifierEdge(tenant_hash)

	err := stmtex.Select(edge_table, edge.ColumnNames(), edge_cond, nil, nil).
		QueryRowsContext(ctx, pub, pub.Dialect())(
		func(scan stmtex.Scanner, _ int) error {
			err := edge.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "failed to scan")
			}

			edge_opt, err := vault.GetChannelNotifierEdge(ctx, pub.DB, pub.Dialect(), edge)
			if err != nil {
				return errors.Wrapf(err, "failed to get a NotifierEdge_option")
			}

			// valied notifier
			err = ValidNotifier(edge_opt)
			if err != nil {
				return errors.Wrapf(err, "valied notifier")
			}

			// notifier factory
			notifier, err := NotifierFactory(channel_uuid, edge_opt)
			if err != nil {
				return errors.Wrapf(err, "channel notifier factory")
			}

			muxer.Notifiers().Add(channel_uuid, notifier)
			return nil
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get notifiers")
	}

	return err
}

func ValidNotifier(edge *channelv3.NotifierEdge_option) (err error) {
	switch edge.NotifierType {
	case channelv3.NotifierTypeConsole:
		err = edge.Console.Valid()
	case channelv3.NotifierTypeRabbitmq:
		err = edge.RabbitMq.Valid()
	case channelv3.NotifierTypeWebhook:
		err = edge.Webhook.Valid()
	case channelv3.NotifierTypeSlackhook:
		err = edge.Slackhook.Valid()
	default:
		err = errors.Errorf("unsupported notifier config%v",
			logs.KVL(
				"opt", edge.NotifierType,
			))
	}

	return
}

func NotifierFactory(uuid string, mc *channelv3.NotifierEdge_option) (notifier Notifier, err error) {
	switch mc.NotifierType {
	case channelv3.NotifierTypeConsole:
		notifier = NewChannelConsole(uuid, mc.Console)
	case channelv3.NotifierTypeRabbitmq:
		notifier = NewChannelRabbitMQ(uuid, mc.RabbitMq)
	case channelv3.NotifierTypeWebhook:
		notifier = NewChannelWebhook(uuid, mc.Webhook)
	case channelv3.NotifierTypeSlackhook:
		notifier = NewChannelSlackhook(uuid, mc.Slackhook)
	default:
		err = errors.Errorf("unsupported notifier option%v",
			logs.KVL(
				"opt", mc.NotifierType,
			))
	}

	return
}
