package managed_channel

import (
	"bytes"
	"strconv"

	"github.com/pkg/errors"
)

type ManagedEventNotifierMux struct {
	// config channelv1.ManagedChannel
	pub Publisher

	notifiers  HashsetNotifier
	formatters HashsetFormatter
}

// var _ ChannelNotifiMuxer = (*ManagedEventNotifierMux)(nil)
// var _ EventNotifierMultiplexer = (*ManagedEventNotifierMux)(nil)

func NewManagedEventNotifierMux( /*cfg channelv1.ManagedChannel*/ ) *ManagedEventNotifierMux {

	mux := &ManagedEventNotifierMux{}
	mux.pub = nil
	mux.notifiers = HashsetNotifier{}
	mux.formatters = HashsetFormatter{}

	// mux.config = cfg //config

	return mux
}

// func (mux ManagedEventNotifierMux) Config() *channelv1.ManagedChannel {
// 	return &mux.config
// }

func (mux ManagedEventNotifierMux) Notifiers() HashsetNotifier {
	return mux.notifiers
}

func (mux ManagedEventNotifierMux) OnNotify(v map[string]interface{}) []error {
	rst := make([]error, 0, len(mux.notifiers))
	for channel_uuid, notifier := range mux.notifiers {
		formatter := mux.formatters[channel_uuid]
		rst = append(rst, notifier.OnNotify(NewMarshalFactory(v, formatter)))
	}
	return rst
}
func (mux ManagedEventNotifierMux) OnNotifyAsync(v map[string]interface{}) []<-chan NotifierFuture {
	futures := make([]<-chan NotifierFuture, 0, len(mux.notifiers))
	for channel_uuid, notifier := range mux.notifiers {
		formatter := mux.formatters[channel_uuid]
		futures = append(futures, OnNotifyAsync(notifier, NewMarshalFactory(v, formatter)))
	}
	return futures
}

// func (mux ManagedEventNotifierMux) Filters() HashsetNotifier {
// 	return mux.notifiers
// }

func (mux ManagedEventNotifierMux) Formatters() HashsetFormatter {
	return mux.formatters
}

func (mux *ManagedEventNotifierMux) Update(v map[string]interface{}) {

	//모든 리스너의 Update 호출 (async)
	futures := mux.OnNotifyAsync(v)

	go func() {
		for _, future := range futures {
			for future := range future {
				if future.Error != nil {
					//업데이트 오류 처리
					mux.EventPublisher().OnNotifierError(
						future.Notifier,
						errors.Wrapf(future.Error, "on notify %v",
							MapString(future.Notifier.Property())))
				}
			}
		}
	}() //!go func()
}

func (mux *ManagedEventNotifierMux) Close() {
}

func (mux *ManagedEventNotifierMux) Regist(pub Publisher) EventNotifierMuxer {
	if pub != nil {
		mux.pub = pub
	}

	if mux.pub != nil {
		mux.pub.SetEventNotifierMuxer(mux)
	}

	return mux
}

func (mux *ManagedEventNotifierMux) EventPublisher() Publisher {
	return mux.pub
}

func MapString(m map[string]string) string {
	buff := bytes.Buffer{}
	for key, value := range m {
		if 0 < buff.Len() {
			buff.WriteString(" ")
		}
		buff.WriteString(key)
		buff.WriteString("=")
		buff.WriteString(strconv.Quote(value))
	}
	return buff.String()
}
