package managed_event

import (
	"fmt"
	"regexp"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

var Invoke func(cluster_uuid, pattern string, i ...interface{}) = func(cluster_uuid, pattern string, i ...interface{}) {}

var _ EventPublisher = (*ManagedEvent)(nil)

type ManagedEvent struct {
	engine *xorm.Engine

	HashsetEventNotifierMultiplexer

	ErrorHandlers         event.HashsetErrorHandlers
	NofitierErrorHandlers HashsetNofitierErrorHandler
}

func NewManagedEvent() *ManagedEvent {

	me := ManagedEvent{}
	me.HashsetEventNotifierMultiplexer = HashsetEventNotifierMultiplexer{}

	me.ErrorHandlers = event.HashsetErrorHandlers{}
	me.NofitierErrorHandlers = HashsetNofitierErrorHandler{}

	return &me
}

func (me ManagedEvent) EventNotifierMultiplexer() HashsetEventNotifierMultiplexer {
	return me.HashsetEventNotifierMultiplexer
}

func (me ManagedEvent) Invoke(cluster_uuid, subscribed_channel string, v ...interface{}) {
	clone := NewManagedEvent()
	clone.engine = me.engine
	clone.ErrorHandlers = me.ErrorHandlers
	clone.NofitierErrorHandlers = me.NofitierErrorHandlers

	//make notifier mux
	if err := clone.BuildNotifierMuxer(cluster_uuid, subscribed_channel); err != nil {
		clone.OnError(errors.Wrapf(err, "build managed event"))
		return
	}

	//update message
	clone.Update(v...)
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

	DefaultErrorHandler_nofitier = func(me *ManagedEvent) func(notifier Notifier, err error) {
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

			record := channelv1.NotifierStatus{}
			record.UuidMeta = metav1.NewUuidMeta()
			record.NotifierType = notifier.Type().String()
			record.NotifierUuid = notifier.Uuid()
			record.Error = fmt.Sprintf("%s%s", err.Error(), stack)

			//이벤트 알림 상태 테이블에 에러 메시지 저장
			if err := vault.NewNotifierStatus(me.Engine().NewSession()).CreateAndRotate(record, globvar.EventNofitierStatusRotateLimit()); err != nil {
				//저장 실패
				me.OnError(errors.Wrapf(err, "save notifier status"))
			}
		}
	}
)

func (me *ManagedEvent) SetEngine(engine *xorm.Engine) *ManagedEvent {
	me.engine = engine

	return me
}

func (me *ManagedEvent) Engine() *xorm.Engine {
	return me.engine
}

func (me *ManagedEvent) Close() {

}

func (me *ManagedEvent) OnError(err error) {
	me.ErrorHandlers.OnError(err)
}
func (me *ManagedEvent) OnNotifierError(notifier Notifier, err error) {
	me.NofitierErrorHandlers.OnError(notifier, err)
}

func (me *ManagedEvent) BuildNotifierMuxer(cluster_uuid, subscribed_channel string) error {
	tx := me.Engine().NewSession()

	//load config
	events, err := vault.NewChannel(tx).Find("cluster_uuid = ?", cluster_uuid)
	if err != nil {
		return errors.Wrapf(err, "find channel by cluster_uuid")
	}
	//subscribed_channel match by regex(event name)
	events_ := make([]channelv1.Channel, 0, len(events))
	for _, event := range events {
		reg, err := regexp.Compile(event.Name)
		if err != nil {
			//regexp compile expr
			return errors.Wrapf(err, "regexp compile%s", logs.KVL(
				"expr", event.Name,
			))
		}

		if ok := reg.MatchString(subscribed_channel); ok {
			events_ = append(events_, event)
		}
	}

	for _, event := range events_ {
		//find edge
		edges, err := vault.NewChannelNotifierEdge(tx).Find("channel_uuid = ?", event.Uuid)
		if err != nil {
			return errors.Wrapf(err, "find channel edge")
		}

		opts := make([]interface{}, 0, 10)
		for _, edge := range edges {
			//get notifier
			opt, err := vault.NewChannelNotifier(tx).Get(edge.NotifierType, edge.NotifierUuid)
			if err != nil {
				return errors.Wrapf(err, "get channel notifier")
			}

			opts = append(opts, opt)
		}
		//new muxer
		muxer := NewManagedEventNotifierMux(event)
		for _, opt := range opts {
			//notifier factory
			notifier, err := NotifierFactory(opt)
			if err != nil {
				return errors.Wrapf(err, "notifier factory")
			}
			//append notifier
			muxer.Notifiers().Add(notifier)
		}

		//regist multiplexer to event publisher
		muxer.Regist(me)
	}

	return nil
}

func NotifierFactory(i interface{}) (new_notifier Notifier, err error) {

	switch opt := i.(type) {
	case *channelv1.NotifierConsole:
		new_notifier = NewConsoleNotifier(opt)
	case *channelv1.NotifierWebhook:
		new_notifier = NewWebhookNotifier(opt)
	case *channelv1.NotifierRabbitMq:
		new_notifier, err = NewRabbitMqNotifier(opt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create rabbitmq notifier%s",
				logs.KVL(
					"opt", opt,
				))
		}
	default:
		return nil, errors.Errorf("unsupported notifier config%s",
			logs.KVL(
				"opt", opt,
			))
	}

	return new_notifier, nil
}
