package managed_event

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
)

type HashsetEventNotifierMultiplexer map[EventNotifierMultiplexer]struct{}

func (hashset HashsetEventNotifierMultiplexer) Add(sub ...EventNotifierMultiplexer) {
	for _, sub := range sub {
		hashset[sub] = struct{}{}
	}
}
func (hashset HashsetEventNotifierMultiplexer) Remove(sub ...EventNotifierMultiplexer) {
	for _, sub := range sub {
		delete(hashset, sub)
	}
}

func (hashset HashsetEventNotifierMultiplexer) Update(v ...interface{}) {
	for mux := range hashset {
		mux.Update(v...)
	}
}

type HashsetNotifier map[Notifier]struct{}

func (hashset HashsetNotifier) Add(notifier ...Notifier) {
	for _, notifier := range notifier {
		hashset[notifier] = struct{}{}
	}
}
func (hashset HashsetNotifier) Remove(notifier ...Notifier) {
	for _, notifier := range notifier {
		delete(hashset, notifier)
	}
}
func (hashset HashsetNotifier) OnNotify(factory MarshalFactoryResult) []error {
	rst := make([]error, 0, len(hashset))
	for notifier := range hashset {
		rst = append(rst, notifier.OnNotify(factory))
	}
	return rst
}
func (hashset HashsetNotifier) OnNotifyAsync(factory MarshalFactoryResult) []<-chan NotifierFuture {
	futures := make([]<-chan NotifierFuture, 0, len(hashset))
	for notifier := range hashset {
		futures = append(futures, OnNotifyAsync(notifier, factory))
	}
	return futures
}

type HashsetNofitierErrorHandler map[uintptr]func(Notifier, error)

func (hashset HashsetNofitierErrorHandler) Add(fn ...func(Notifier, error)) HashsetNofitierErrorHandler {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		hashset[ptr] = fn
	}

	return hashset
}
func (hashset HashsetNofitierErrorHandler) Remove(fn ...func(error)) HashsetNofitierErrorHandler {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		delete(hashset, ptr)
	}

	return hashset
}
func (hashset HashsetNofitierErrorHandler) OnError(notifier Notifier, err error) {
	for _, handler := range hashset {
		handler(notifier, err)
	}
}

type EventNotifiMuxConfigHolder interface {
	Config() *channelv1.Channel
}

type NotifierFuture struct {
	Notifier Notifier
	Error    error
}

//Notifier
type Notifier interface {
	Type() fmt.Stringer          //리스너 타입
	Uuid() string                //uuid
	Property() map[string]string //요약
	// PropertyString() string                                   //요약
	OnNotify(MarshalFactoryResult) error //알림 발생
	Close()                              //리스너 종료
}

func OnNotifyAsync(notifier Notifier, factory MarshalFactoryResult) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{Notifier: notifier, Error: notifier.OnNotify(factory)}
	}()

	return future
}

//EventNotifierMultiplexer
type EventNotifierMultiplexer interface {
	Notifiers() HashsetNotifier //Notifiers
	Update(...interface{})      //Update 발생
	EventPublisher() EventPublisher
	Regist(EventPublisher) EventNotifierMultiplexer
	Close() //이벤트 구독 취소 //전체 Notifier 제거
}

//EventPublisher
type EventPublisher interface {
	EventNotifierMultiplexer() HashsetEventNotifierMultiplexer
	Close()
	OnError(error)
	OnNotifierError(Notifier, error)
}
type MarshalFactoryResult func(string) ([][]byte, error)

func NewMarshalFactory(v ...interface{}) func(string) ([][]byte, error) {
	mux := sync.Mutex{}
	m := make(map[string][][]byte)

	return func(mime string) ([][]byte, error) {
		mux.Lock()
		defer mux.Unlock()

		//이미 저장된 데이터가 있으면 저장된 데이터를 리턴
		if _, ok := m[mime]; ok {
			return m[mime], nil
		}

		//저장된 데이터가 없으면 데이터 만들어서 저장
		m[mime] = make([][]byte, len(v))
		for i, v := range v {
			var (
				b   []byte
				err error
			)
			switch strings.ToLower(mime) {
			case "application/json":
				b, err = json.Marshal(v)

			case "application/xml":
				b, err = xml.Marshal(v)
			default:
				err = errors.Errorf("unsupported Content-Type")
			}

			if err != nil {
				return nil, errors.Wrapf(err, "marshal factory%s",
					logs.KVL(
						"item", v,
						"mime", mime,
					))
			}
			m[mime][i] = b
		}

		return m[mime], nil
	}
}
