package managed_event

import (
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
)

type ManagedEventNotifierMux struct {
	config eventv1.Event
	pub    EventPublisher

	notifiers HashsetNotifier
}

var _ EventNotifiMuxConfigHolder = (*ManagedEventNotifierMux)(nil)
var _ EventNotifierMultiplexer = (*ManagedEventNotifierMux)(nil)

func NewManagedEventNotifierMux(cfg eventv1.Event) *ManagedEventNotifierMux {

	mux := &ManagedEventNotifierMux{}
	mux.notifiers = HashsetNotifier{}

	mux.config = cfg //config

	return mux
}

func (mux ManagedEventNotifierMux) Config() *eventv1.Event {
	return &mux.config
}

func (mux ManagedEventNotifierMux) Notifiers() HashsetNotifier {
	return mux.notifiers
}

func (mux *ManagedEventNotifierMux) Update(v ...interface{}) {
	factory := event.NewMarshalFactory(v...)
	//모든 리스너의 Update 호출 (async)
	futures := mux.notifiers.OnNotifyAsync(factory)

	go func() {
		for _, future := range futures {
			for future := range future {
				if future.Error != nil {
					//업데이트 오류 처리
					mux.EventPublisher().OnNotifierError(
						future.Notifier,
						errors.Wrapf(future.Error, "event notify%s %s",
							logs.KVL(
								"cluster_uuid", mux.Config().ClusterUuid,
								"event_uuid", mux.Config().Uuid,
								"event_pattern", mux.Config().Pattern,
							),
							event.MapString(future.Notifier.Property())))

				}
			}
		}
	}() //!go func()
}

func (mux *ManagedEventNotifierMux) Close() {
}

func (mux *ManagedEventNotifierMux) Regist(pub EventPublisher) EventNotifierMultiplexer {
	if pub != nil {
		mux.pub = pub
	}

	if mux.pub != nil {
		mux.pub.EventNotifierMultiplexer().Add(mux)
	}

	return mux
}

func (mux *ManagedEventNotifierMux) EventPublisher() EventPublisher {
	return mux.pub
}
