package event

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type HashsetEventNotifierMux map[EventNotifierMultiplexer]struct{}

func (hashset HashsetEventNotifierMux) Add(sub ...EventNotifierMultiplexer) {
	for _, sub := range sub {
		hashset[sub] = struct{}{}
	}
}
func (hashset HashsetEventNotifierMux) Remove(sub ...EventNotifierMultiplexer) {
	for _, sub := range sub {
		delete(hashset, sub)
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

type NotifierFuture struct {
	Notifier Notifier
	Error    error
}

//Notifier
type Notifier interface {
	Type() fmt.Stringer          //리스너 타입
	Property() map[string]string //요약
	// PropertyString() string                                   //요약
	OnNotify(MarshalFactoryResult) error //알림 발생
	// OnNotifyAsync(MarshalFactoryResult) <-chan NotifierFuture //알림 발생 (async)
	Regist(EventNotifierMultiplexer) //이벤트 구독
	Close()                          //리스너 종료
}

func OnNotifyAsync(notifier Notifier, factory MarshalFactoryResult) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{Notifier: notifier, Error: notifier.OnNotify(factory)}
	}()

	return future
}

type EventNotifiMuxConfigHolder interface {
	Config() *EventNotifierMuxerConfig //설정
}

//EventNotifierMultiplexer
type EventNotifierMultiplexer interface {
	Notifiers() HashsetNotifier    //Notifiers
	Update(string, ...interface{}) //Update 발생
	Regist(EventPublisher)         //EventPublisher 등록
	Close()                        //이벤트 구독 취소 //전체 Notifier 제거
	// ErrorHandlers() HashsetErrorHandlers //에러 핸들러
}

//EventPublisher
type EventPublisher interface {
	Publish(string, ...interface{})
	NotifierMuxers() HashsetEventNotifierMux //Notifiers
	Close()                                  //이벤트
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
			switch mime {
			case "application/json":
				b, err := json.Marshal(v)
				if err != nil {
					return nil, errors.Wrapf(err, "json marshal%s",
						logs.KVL(
							"item", v,
						))
				}
				m[mime][i] = b
			case "application/xml":

				switch v := v.(type) {
				case map[string]interface{}:

					sm := StringMap{}
					sm.Map = map[string]interface{}{}
					sm.EntryName = "root"
					if v["event_name"] != nil {
						if _, ok := v["event_name"].(string); ok {
							sm.EntryName = v["event_name"].(string)
						}
					}
					for k, v := range v {
						sm.Map[k] = v
					}

					b, err := xml.Marshal(sm)
					if err != nil {
						return nil, errors.Wrapf(err, "json marshal%s",
							logs.KVL(
								"item", v,
							))
					}
					m[mime][i] = b
				}

			default:
				return nil, fmt.Errorf("unsupported Content-Type=%s", strconv.Quote(mime))
			}
		}

		return m[mime], nil
	}
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

type StringMap struct {
	EntryName string `xml:"-"`
	Map       map[string]interface{}
}

func (s StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = s.EntryName
	tokens := []xml.Token{start}
	for k, v := range s.Map {
		v := fmt.Sprintf("%v", v)
		t := xml.StartElement{Name: xml.Name{"", k}}
		tokens = append(tokens, t, xml.CharData(v), xml.EndElement{t.Name})
	}
	tokens = append(tokens, xml.EndElement{start.Name})
	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}
