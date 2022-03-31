package event

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type HashsetEventSubscribers map[EventSubscriber]struct{}

func (hashset HashsetEventSubscribers) Add(sub ...EventSubscriber) {
	for _, sub := range sub {
		hashset[sub] = struct{}{}
	}
}
func (hashset HashsetEventSubscribers) Remove(sub ...EventSubscriber) {
	for _, sub := range sub {
		delete(hashset, sub)
	}
}

type HashsetNotifiers map[Notifier]struct{}

func (hashset HashsetNotifiers) Add(notifier ...Notifier) {
	for _, notifier := range notifier {
		hashset[notifier] = struct{}{}
	}
}
func (hashset HashsetNotifiers) Remove(notifier ...Notifier) {
	for _, notifier := range notifier {
		delete(hashset, notifier)
	}
}
func (hashset HashsetNotifiers) OnNotify(factory MarshalFactory) []error {
	rst := make([]error, 0, len(hashset))
	for notifier := range hashset {
		rst = append(rst, notifier.OnNotify(factory))
	}
	return rst
}
func (hashset HashsetNotifiers) OnNotifyAsync(factory MarshalFactory) []<-chan NotifierFuture {
	futures := make([]<-chan NotifierFuture, 0, len(hashset))
	for notifier := range hashset {
		futures = append(futures, notifier.OnNotifyAsync(factory))
	}
	return futures
}

type HashsetErrorHandlers map[uintptr]func(error)

func (hashset HashsetErrorHandlers) Add(fn ...func(error)) {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		hashset[ptr] = fn
	}
}
func (hashset HashsetErrorHandlers) Remove(fn ...func(error)) {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		delete(hashset, ptr)
	}
}
func (hashset HashsetErrorHandlers) OnError(err error) {
	for _, handler := range hashset {
		handler(err)
	}
}

type NotifierFuture struct {
	Notifier Notifier
	Error    error
}

//Notifier
type Notifier interface {
	Type() string                                       //리스너 타입
	Property() map[string]string                        //요약
	PropertyString() string                             //요약
	OnNotify(MarshalFactory) error                      //알림 발생
	OnNotifyAsync(MarshalFactory) <-chan NotifierFuture //알림 발생 (async)
	Regist(EventSubscriber)                             //이벤트 구독
	Close()                                             //리스너 종료
}

//EventSubscriber
type EventSubscriber interface {
	Config() *EventSubscribeConfig //설정
	Notifiers() HashsetNotifiers   //Notifiers
	Update(string, ...interface{}) //Update 발생
	Regist(EventPublisher)         //EventPublisher 등록
	Close()                        //이벤트 구독 취소 //전체 Notifier 제거
	// ErrorHandlers() HashsetErrorHandlers //에러 핸들러
}

//EventPublisher
type EventPublisher interface {
	Publish(string, ...interface{})
	Subscribers() HashsetEventSubscribers //Notifiers
	Close()                               //이벤트
}

type MarshalFactory func(string) ([][]byte, error)

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
		m[mime] = make([][]byte, 0, len(v))
		for _, v := range v {
			var (
				b   []byte
				err error
			)
			switch mime {
			case "application/json":
				if b, err = json.Marshal(v); err != nil {
					return nil, errors.Wrapf(err, "json marshal%s",
						logs.KVL(
							"item", v,
						))
				}
			default:
				return nil, fmt.Errorf("unsupported Content-Type=%s", strconv.Quote(mime))
			}

			m[mime] = append(m[mime], b)
		}

		return m[mime], nil
	}
}
