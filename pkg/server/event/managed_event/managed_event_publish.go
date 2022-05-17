package managed_event

import (
	"fmt"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

var Invoke func(cluster_uuid, pattern string, i ...interface{}) = func(cluster_uuid, pattern string, i ...interface{}) {}

var _ EventPublisher = (*ManagedEvent)(nil)

type ManagedEvent struct {
	engine *xorm.Engine

	ErrorHandlers         event.HashsetErrorHandlers
	NofitierErrorHandlers HashsetNofitierErrorHandler
}

func NewManagedEvent() *ManagedEvent {

	me := ManagedEvent{}
	me.ErrorHandlers = event.HashsetErrorHandlers{}
	me.NofitierErrorHandlers = HashsetNofitierErrorHandler{}

	return &me
}

func (me ManagedEvent) Invoke(cluster_uuid, pattern string, i ...interface{}) {

	//make notifier mux
	mux, err := me.BuildNotifierMuxer(cluster_uuid, pattern)
	if err != nil {
		me.OnError(errors.Wrapf(err, "build managed event"))
		return
	}

	//update message
	if mux != nil {
		mux.Update(pattern, i...)
	}
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
					me.OnError(errors.Wrapf(err, "recover notifier%s",
						event.MapString(notifier.Property())))
				} else {
					me.OnError(errors.Errorf("recover notifier %s: %+v",
						event.MapString(notifier.Property()), r))
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

			record := eventv1.EventNotifierStatus{}
			record.UuidMeta = metav1.NewUuidMeta()
			record.NotifierType = notifier.Type().String()
			record.NotifierUuid = notifier.Uuid()
			record.Error = fmt.Sprintf("%s%s", err.Error(), stack)

			//이벤트 알림 상태 테이블에 에러 메시지 저장
			if err := vault.NewEventNotifierStatus(me.Engine().NewSession()).CreateAndRotate(record, globvar.EventNofitierStatusRotateLimit()); err != nil {
				//저장 실패
				me.OnError(errors.Wrapf(err, "notifier%s",
					event.MapString(notifier.Property())))
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

func (me *ManagedEvent) BuildNotifierMuxer(cluster_uuid, pattern string) (EventNotifierMuxer, error) {

	tx := me.Engine().NewSession()

	var mux EventNotifierMuxer

	//load config
	events, err := vault.NewEvent(tx).Find("cluster_uuid = ? AND pattern = ?", cluster_uuid, pattern)
	if err != nil {
		return mux, errors.Wrapf(err, "find event")
	}

	for _, event := range events {
		//new mux
		mux = NewManagedEventNotifierMux(event)

		//get edge
		edges, err := vault.NewEventNotifierEdge(tx).Find("event_uuid = ?", event.Uuid)
		if err != nil {
			return mux, errors.Wrapf(err, "find event")
		}

		for _, edge := range edges {
			//get notifier
			i, err := vault.NewEventNotifier(tx).Get(edge.NotifierType, edge.NotifierUuid)
			if err != nil {
				return mux, errors.Wrapf(err, "find event")
			}

			if _, ok := i[edge.NotifierType]; !ok {
				continue
			}

			//notifier factory
			notifier, err := NotifierFactory(i[edge.NotifierType])
			if err != nil {
				return mux, errors.Wrapf(err, "notifier factory")
			}
			//append notifier
			mux.Notifiers().Add(notifier)
		}

		mux.Regist(me)
	}

	return mux, nil
}

func NotifierFactory(i interface{}) (new_notifier Notifier, err error) {

	switch opt := i.(type) {
	case *eventv1.EventNotifierConsole:
		new_notifier = NewConsoleNotifier(opt)
	case *eventv1.EventNotifierWebhook:
		new_notifier = NewWebhookNotifier(opt)
	case *eventv1.EventNotifierRabbitMq:
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
