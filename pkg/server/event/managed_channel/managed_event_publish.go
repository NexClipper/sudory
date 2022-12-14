package managed_channel

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/NexClipper/logger"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/pkg/errors"
)

var InvokeByChannelUuid func(tenant_hash string, channel_uuid string, v []map[string]interface{}) = func(tenant_hash string, channel_uuid string, v []map[string]interface{}) {}

var InvokeByEventCategory func(tenant_hash string, ec channelv3.EventCategory, v []map[string]interface{}) = func(tenant_hash string, ec channelv3.EventCategory, v []map[string]interface{}) {}

// var _ EventPublisher = (*ManagedChannel)(nil)

type HashsetErrorHandlers map[uintptr]func(error)

func (hashset HashsetErrorHandlers) Add(fn ...func(error)) HashsetErrorHandlers {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		hashset[ptr] = fn
	}

	return hashset
}
func (hashset HashsetErrorHandlers) Remove(fn ...func(error)) HashsetErrorHandlers {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		delete(hashset, ptr)
	}

	return hashset
}
func (hashset HashsetErrorHandlers) OnError(err error) {
	for _, handler := range hashset {
		handler(err)
	}
}

type Event struct {
	ctx context.Context
	*sql.DB
	dialect excute.SqlExcutor
	// *vanilla.SqlDbEx

	// HashsetEventNotifierMuxer
	EventNotifierMuxer

	ErrorHandlers         HashsetErrorHandlers
	NofitierErrorHandlers HashsetNofitierErrorHandler
}

func NewEvent(db *sql.DB, dialect excute.SqlExcutor) *Event {

	me := Event{}
	me.ctx = context.TODO()
	me.DB = db
	me.dialect = dialect
	// me.SqlDbEx = &vanilla.SqlDbEx{DB: db}
	// me.HashsetEventNotifierMuxer = HashsetEventNotifierMuxer{}

	me.ErrorHandlers = HashsetErrorHandlers{}
	me.NofitierErrorHandlers = HashsetNofitierErrorHandler{}

	return &me
}

// func (pub *Event) Dialect() string {
// 	return pub.dialect
// }

func (pub *Event) SetEventNotifierMuxer(mux EventNotifierMuxer) {
	pub.EventNotifierMuxer = mux
}

// func (pub Event) EventNotifierMuxer() HashsetEventNotifierMuxer {
// 	return pub.HashsetEventNotifierMuxer
// }

func (pub Event) InvokeByChannelUuid(tenant_hash string, channel_uuid string, v []map[string]interface{}) {
	clone := NewEvent(pub.DB, pub.dialect)
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

	// notifier close
	clone.Close()
}

func (pub Event) BuildChannelFormatter(channel_uuid string) (err error) {

	eq_uuid := stmt.Equal("uuid", channel_uuid)
	format := channelv3.Format{}
	err = pub.dialect.QueryRows(format.TableName(), format.ColumnNames(), eq_uuid, nil, nil)(pub.ctx, pub)(
		func(scan excute.Scanner, _ int) error {
			err := format.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
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

			return err
		})
	err = errors.Wrapf(err, "build channel formatter")

	return
}

func (pub Event) InvokeByEventCategory(tenant_hash string, ec channelv3.EventCategory, v []map[string]interface{}) {
	clone := NewEvent(pub.DB, pub.dialect)
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

	// notifier close
	clone.Close()
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

			if err_ := vault.CreateChannelStatus(pub.ctx, pub.DB, pub.dialect, uuid, message, created, globvar.Event.NofitierStatusRotateLimit()); err_ != nil {
				err_ = errors.Wrapf(err_, "failed to logging to channel status")
				pub.ErrorHandlers.OnError(err_)
			}
		}
	}
)

func (pub *Event) Close() {
	for _, notifier := range pub.EventNotifierMuxer.Notifiers() {
		notifier.Close()
	}
}

func (pub *Event) OnError(err error) {
	pub.ErrorHandlers.OnError(err)
}
func (pub *Event) OnNotifierError(notifier Notifier, err error) {
	pub.NofitierErrorHandlers.OnError(notifier, err)
}

func (pub *Event) BuildMuxerByEventCategory(tenant_hash string, event_category channelv3.EventCategory) (err error) {

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
	var set_channel_uuid = map[string]struct{}{}
	channel_table := channelv3.TableNameWithTenant_ManagedChannel(tenant_hash)
	err = pub.dialect.QueryRows(channel_table, column_names, channel_cond, nil, nil)(pub.ctx, pub)(
		func(scan excute.Scanner, _ int) error {

			var channel_uuid string
			err := scan.Scan(&channel_uuid)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			set_channel_uuid[channel_uuid] = struct{}{}

			return err
		})
	if err != nil {
		return
	}

	for channel_uuid := range set_channel_uuid {
		var set_edge = map[channelv3.NotifierEdge]struct{}{}
		// get notifier edge with option
		var edge channelv3.NotifierEdge
		edge.Uuid = channel_uuid
		edge_cond := stmt.Equal("uuid", edge.Uuid)

		err = pub.dialect.QueryRows(edge.TableName(), edge.ColumnNames(), edge_cond, nil, nil)(pub.ctx, pub)(
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
			return
		}

		for edge := range set_edge {
			edge_opt, err := vault.GetChannelNotifierEdge(pub.ctx, pub.DB, pub.dialect, edge)
			if err != nil {
				return errors.Wrapf(err, "failed to get a NotifierEdge_option")
			}

			// valid notifier
			err = ValidNotifier(edge_opt)
			if err != nil {
				return errors.Wrapf(err, "valid notifier")
			}

			// notifier factory
			notifier, err := NotifierFactory(channel_uuid, edge_opt)
			if err != nil {
				return errors.Wrapf(err, "channel notifier factory")
			}

			// append notifier
			muxer.Notifiers().Add(channel_uuid, notifier)
		}
	}
	return nil
}

func (pub *Event) BuildMuxerByChannelUuid(tenant_hash string, channel_uuid string) error {

	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(pub)

	// get notifier edge with option
	var set_edge = map[channelv3.NotifierEdge]struct{}{}
	var edge channelv3.NotifierEdge
	edge.Uuid = channel_uuid
	edge_cond := stmt.Equal("uuid", channel_uuid)
	edge_table := channelv3.TableNameWithTenant_NotifierEdge(tenant_hash)

	err := pub.dialect.QueryRows(edge_table, edge.ColumnNames(), edge_cond, nil, nil)(pub.ctx, pub)(
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
		return errors.Wrapf(err, "failed to get notifiers")
	}

	for edge := range set_edge {
		edge_opt, err := vault.GetChannelNotifierEdge(pub.ctx, pub.DB, pub.dialect, edge)
		if err != nil {
			return errors.Wrapf(err, "failed to get a NotifierEdge_option")
		}

		// valid notifier
		err = ValidNotifier(edge_opt)
		if err != nil {
			return errors.Wrapf(err, "valid notifier")
		}

		// notifier factory
		notifier, err := NotifierFactory(channel_uuid, edge_opt)
		if err != nil {
			return errors.Wrapf(err, "channel notifier factory")
		}

		muxer.Notifiers().Add(channel_uuid, notifier)
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
