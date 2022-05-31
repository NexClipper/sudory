package managed_event

import (
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
)

type ManagedEventNotifierMux struct {
	config channelv1.Channel
	pub    EventPublisher

	notifiers HashsetNotifier
}

var _ EventNotifiMuxConfigHolder = (*ManagedEventNotifierMux)(nil)
var _ EventNotifierMultiplexer = (*ManagedEventNotifierMux)(nil)

func NewManagedEventNotifierMux(cfg channelv1.Channel) *ManagedEventNotifierMux {

	mux := &ManagedEventNotifierMux{}
	mux.notifiers = HashsetNotifier{}

	mux.config = cfg //config

	return mux
}

func (mux ManagedEventNotifierMux) Config() *channelv1.Channel {
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
						errors.Wrapf(future.Error, "on notify%s %s",
							logs.KVL(
								"cluster_uuid", mux.Config().ClusterUuid,
								"channel_uuid", mux.Config().Uuid,
								"channel_name", mux.Config().Name,
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
