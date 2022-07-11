package managed_channel

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/NexClipper/logger"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/pkg/errors"
)

var InvokeByChannelUuid func(channel_uuid string, v map[string]interface{}) = func(channel_uuid string, v map[string]interface{}) {}

var InvokeByEventCategory func(ec channelv2.EventCategory, v map[string]interface{}) = func(ec channelv2.EventCategory, v map[string]interface{}) {}

// var _ EventPublisher = (*ManagedChannel)(nil)

type Event struct {
	// db *sql.DB
	*vanilla.SqlDbEx

	// HashsetEventNotifierMuxer
	EventNotifierMuxer

	ErrorHandlers         event.HashsetErrorHandlers
	NofitierErrorHandlers HashsetNofitierErrorHandler
}

func NewEvent(db *sql.DB) *Event {

	me := Event{}
	me.SqlDbEx = vanilla.NewSqlDbEx(db)
	// me.HashsetEventNotifierMuxer = HashsetEventNotifierMuxer{}

	me.ErrorHandlers = event.HashsetErrorHandlers{}
	me.NofitierErrorHandlers = HashsetNofitierErrorHandler{}

	return &me
}

func (pub *Event) SetEventNotifierMuxer(mux EventNotifierMuxer) {
	pub.EventNotifierMuxer = mux
}

// func (pub Event) EventNotifierMuxer() HashsetEventNotifierMuxer {
// 	return pub.HashsetEventNotifierMuxer
// }

func (pub Event) InvokeByChannelUuid(channel_uuid string, v map[string]interface{}) {
	clone := NewEvent(pub.DB)
	clone.ErrorHandlers = pub.ErrorHandlers
	clone.NofitierErrorHandlers = pub.NofitierErrorHandlers

	//build channel muxer
	if err := clone.BuildMuxerByChannelUuid(channel_uuid); err != nil {
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
	eq_uuid := vanilla.Equal("uuid", channel_uuid)
	format := channelv2.Format{}
	err = vanilla.Stmt.Select(format.TableName(), format.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRows(pub)(func(scan vanilla.Scanner, _ int) (err error) {
		err = format.Scan(scan)
		if err != nil {
			return
		}
		switch format.FormatType {
		case channelv2.FormatTypeFields:
			formatter := &Formatter_fields{
				FormatData: format.FormatData,
			}
			pub.EventNotifierMuxer.Formatters().Add(channel_uuid, formatter)
		case channelv2.FormatTypeJq:
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

func (pub Event) InvokeByEventCategory(ec channelv2.EventCategory, v map[string]interface{}) {
	clone := NewEvent(pub.DB)
	clone.ErrorHandlers = pub.ErrorHandlers
	clone.NofitierErrorHandlers = pub.NofitierErrorHandlers

	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(clone)

	//build channel muxer
	if err := clone.BuildMuxerByEventCategory(ec); err != nil {
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

	DefaultErrorHandler_notifier = func(me *Event) func(notifier Notifier, err error) {
		return func(notifier Notifier, err error) {
			defer func() {
				r := recover()

				if r == nil {
					return
				}

				if err, ok := r.(error); ok {
					me.OnError(errors.Wrapf(err, "recover notifier error handler"))
				} else {
					me.OnError(errors.Errorf("notifier error handler recover='%+v'", r))
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

			if err_ := vault.CreateChannelStatus(me.DB, uuid, message, created, globvar.EventNofitierStatusRotateLimit()); err_ != nil {
				err_ = errors.Wrapf(err_, "failed to logging to channel status")
				me.ErrorHandlers.OnError(err_)
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

func (pub *Event) BuildMuxerByEventCategory(event_category channelv2.EventCategory) (err error) {
	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(pub)

	find_channel_cond := vanilla.And(
		vanilla.IsNull("deleted"),
		vanilla.Equal("event_category", int(event_category)),
	).Parse()
	column_names := []string{"uuid"}

	find_channel := channelv2.ManagedChannel{}
	err = vanilla.Stmt.Select(find_channel.TableName(), column_names, find_channel_cond, nil, nil).
		QueryRows(pub)(func(scan vanilla.Scanner, _ int) (err error) {

		var channel_uuid string
		err = scan.Scan(&channel_uuid)
		if err != nil {
			return
		}

		// get notifier edge with option
		eq_uuid := vanilla.Equal("uuid", channel_uuid)

		notifier_edge_option := new(channelv2.NotifierEdge_option)
		err = vanilla.Stmt.Select(notifier_edge_option.TableName(), notifier_edge_option.ColumnNames(), eq_uuid.Parse(), nil, nil).
			QueryRows(pub)(func(scan vanilla.Scanner, _ int) (err error) {
			err = notifier_edge_option.Scan(scan)
			return
		})
		err = errors.Wrapf(err, "failed to query from NotifierEdge_option")
		if err != nil {
			return
		}

		// valied notifier
		err = ValidNotifier(notifier_edge_option)
		err = errors.Wrapf(err, "valied notifier")
		if err != nil {
			return
		}

		// notifier factory
		notifier, err := NotifierFactory(channel_uuid, notifier_edge_option)
		err = errors.Wrapf(err, "channel notifier factory")
		if err != nil {
			return
		}

		// append notifier
		muxer.Notifiers().Add(channel_uuid, notifier)

		return
	})
	if err != nil {
		return
	}

	return nil
}

func (pub *Event) BuildMuxerByChannelUuid(channel_uuid string) (err error) {

	// new muxer
	muxer := NewManagedEventNotifierMux()
	// regist muxer to event publisher
	muxer.Regist(pub)

	// get notifier edge with option
	eq_uuid := vanilla.Equal("uuid", channel_uuid)

	notifier_edge_option := new(channelv2.NotifierEdge_option)
	err = vanilla.Stmt.Select(notifier_edge_option.TableName(), notifier_edge_option.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRows(pub)(func(scan vanilla.Scanner, _ int) (err error) {
		err = notifier_edge_option.Scan(scan)
		return
	})
	err = errors.Wrapf(err, "failed to query from NotifierEdge_option")
	if err != nil {
		return
	}

	// valied notifier
	err = ValidNotifier(notifier_edge_option)
	err = errors.Wrapf(err, "valied notifier")
	if err != nil {
		return
	}

	// notifier factory
	notifier, err := NotifierFactory(channel_uuid, notifier_edge_option)
	err = errors.Wrapf(err, "channel notifier factory")
	if err != nil {
		return
	}

	muxer.Notifiers().Add(channel_uuid, notifier)

	return
}

func ValidNotifier(mc *channelv2.NotifierEdge_option) (err error) {
	switch mc.NotifierType {
	case channelv2.NotifierTypeConsole:

	case channelv2.NotifierTypeRabbitmq:
		if len(mc.RabbitMq.Url) == 0 {
			err = errors.Errorf("missing url")
		}
		if len(mc.RabbitMq.ChannelPublish.Exchange.String) == 0 &&
			len(mc.RabbitMq.ChannelPublish.RoutingKey.String) == 0 {
			err = errors.Errorf("missing exchange or routing-key")
		}
	case channelv2.NotifierTypeWebhook:
		if len(mc.Webhook.Method) == 0 {
			err = errors.Errorf("missing method")
		}
		if len(mc.Webhook.Url) == 0 {
			err = errors.Errorf("missing url")
		}
	default:
		err = errors.Errorf("unsupported notifier config%v",
			logs.KVL(
				"opt", mc.NotifierType,
			))
	}

	return
}

func NotifierFactory(uuid string, mc *channelv2.NotifierEdge_option) (notifier Notifier, err error) {
	switch mc.NotifierType {
	case channelv2.NotifierTypeConsole:
		notifier = NewChannelConsole(uuid, mc.Console)
	case channelv2.NotifierTypeRabbitmq:
		notifier = NewChannelRabbitMQ(uuid, mc.RabbitMq)
	case channelv2.NotifierTypeWebhook:
		notifier = NewChannelWebhook(uuid, mc.Webhook)
	default:
		err = errors.Errorf("unsupported notifier option%v",
			logs.KVL(
				"opt", mc.NotifierType,
			))
	}

	return
}
