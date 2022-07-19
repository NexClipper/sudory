package managed_channel

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type HashsetFormatter map[string]Formatter

func (hashset HashsetFormatter) Add(channel_uuid string, formatter Formatter) {
	hashset[channel_uuid] = formatter
}
func (hashset HashsetFormatter) Remove(channel_uuid string) {
	delete(hashset, channel_uuid)
}

// type HashsetEventNotifierMuxer map[EventNotifierMuxer]struct{}

// func (hashset HashsetEventNotifierMuxer) Add(sub ...EventNotifierMuxer) {
// 	for _, sub := range sub {
// 		hashset[sub] = struct{}{}
// 	}
// }
// func (hashset HashsetEventNotifierMuxer) Remove(sub ...EventNotifierMuxer) {
// 	for _, sub := range sub {
// 		delete(hashset, sub)
// 	}
// }

// func (hashset HashsetEventNotifierMuxer) Update(v interface{}) {
// 	for mux := range hashset {
// 		mux.Update(v)
// 	}
// }

type HashsetNotifier map[string]Notifier

func (hashset HashsetNotifier) Add(channel_uuid string, notifier Notifier) {
	hashset[channel_uuid] = notifier
}
func (hashset HashsetNotifier) Remove(channel_uuid string) {
	delete(hashset, channel_uuid)
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

// type ChannelNotifiMuxer interface {
// 	Config() *channelv1.ManagedChannel
// }

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
	OnNotify(*MarshalFactory) error //알림 발생
	Close()                         //리스너 종료
}

func OnNotifyAsync(notifier Notifier, factory *MarshalFactory) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{Notifier: notifier, Error: notifier.OnNotify(factory)}
	}()

	return future
}

//EventNotifierMuxer
type EventNotifierMuxer interface {
	Notifiers() HashsetNotifier    // Notifiers
	Formatters() HashsetFormatter  // Formatter
	Update(map[string]interface{}) // Update 발생
	EventPublisher() Publisher
	Regist(Publisher) EventNotifierMuxer
	Close() // 이벤트 구독 취소 // 전체 Notifier 제거
}

//Publisher
type Publisher interface {
	SetEventNotifierMuxer(EventNotifierMuxer)
	Close()
	OnError(error)
	OnNotifierError(Notifier, error)
}

type MarshalFactory struct {
	mux   sync.Mutex
	cache map[string][]byte

	Value     map[string]interface{}
	Formatter Formatter

	format_once   sync.Once
	formated_data interface{}
}

func NewMarshalFactory(v map[string]interface{}, format Formatter) *MarshalFactory {
	return &MarshalFactory{
		Value:     v,
		Formatter: format,
	}
}

func (ms *MarshalFactory) Marshal(mime string) (b []byte, err error) {
	ms.mux.Lock()
	defer ms.mux.Unlock()

	if ms.cache == nil {
		ms.cache = make(map[string][]byte)
	}

	// 이미 저장된 데이터가 있으면 저장된 데이터를 리턴
	if _, ok := ms.cache[mime]; ok {
		return ms.cache[mime], nil // hit
	}

	ms.format_once.Do(func() {
		if ms.Formatter != nil {
			ms.formated_data, err = ms.Formatter.Format(ms.Value)
			err = errors.Wrapf(err, "format%v",
				logs.KVL(
					"format_type", ms.Formatter.Type(),
				))
		} else {
			ms.formated_data = ms.Value
		}
	})

	if err != nil {
		return
	}

	//저장된 데이터가 없으면 데이터 만들어서 저장
	// m[mime] = make([]byte, len(v))
	// for i, v := range v {

	switch strings.ToLower(mime) {
	case "application/json":
		b, err = json.Marshal(ms.formated_data)
	case "application/xml":
		b, err = xml.Marshal(ms.formated_data)
	case "application/x-www-form-urlencoded":
		b, err = urlencoded(ms.formated_data)
	case "text/plain":
		b, err = text_plain(ms.formated_data)
	default:
		err = errors.Errorf("unsupported Content-Type")
	}
	err = errors.Wrapf(err, "marshal factory%s",
		logs.KVL(
			"item", ms.formated_data,
			"mime", mime,
		))
	if err != nil {
		return
	}
	ms.cache[mime] = b
	// }

	return ms.cache[mime], nil
}

// type MarshalFactoryResult func(string) ([]byte, error)

// func NewMarshalFactory(v map[string]interface{}, formatter Formatter) func(string) ([]byte, error) {
// 	mux := sync.Mutex{}
// 	m := make(map[string][]byte)

// 	var err error
// 	var v_formated interface{}
// 	if formatter != nil {
// 		v_formated, err = formatter.Format(v)
// 	} else {
// 		v_formated = v
// 	}

// 	if err != nil {
// 		return func(mime string) ([]byte, error) {
// 			return nil, err
// 		}
// 	}

// 	return func(mime string) ([]byte, error) {
// 		mux.Lock()
// 		defer mux.Unlock()

// 		//이미 저장된 데이터가 있으면 저장된 데이터를 리턴
// 		if _, ok := m[mime]; ok {
// 			return m[mime], nil
// 		}

// 		//저장된 데이터가 없으면 데이터 만들어서 저장
// 		// m[mime] = make([]byte, len(v))
// 		// for i, v := range v {
// 		var (
// 			b   []byte
// 			err error
// 		)
// 		switch strings.ToLower(mime) {
// 		case "application/json":
// 			b, err = json.Marshal(v_formated)

// 		case "application/xml":
// 			b, err = xml.Marshal(v_formated)
// 		case "application/x-www-form-urlencoded":
// 			b, err = urlencoded(v_formated)
// 		case "text/plain":
// 			b, err = text_plain(v_formated)
// 		default:
// 			err = errors.Errorf("unsupported Content-Type")
// 		}

// 		if err != nil {
// 			return nil, errors.Wrapf(err, "marshal factory%s",
// 				logs.KVL(
// 					"item", v_formated,
// 					"mime", mime,
// 				))
// 		}
// 		m[mime] = b
// 		// }

// 		return m[mime], nil
// 	}
// }

// application/x-www-form-urlencoded
func urlencoded(v interface{}) ([]byte, error) {
	conv_map := func(v map[string]interface{}) string {
		s := make([]string, 0, len(v))
		for k, v := range v {
			s = append(s, fmt.Sprintf("%v=%v", k, v))
		}
		return strings.Join(s, "&")
	}

	conv_keyval := func(v map[string]string) string {
		s := make([]string, 0, len(v))
		for k, v := range v {
			s = append(s, fmt.Sprintf("%v=%v", k, v))
		}
		return strings.Join(s, "&")
	}

	switch v := v.(type) {
	case map[string]interface{}:
		return []byte(conv_map(v)), nil
	case map[string]string:
		return []byte(conv_keyval(v)), nil
	default:
		return []byte{}, fmt.Errorf("this type is not supported")
	}
}

// text/plain
//  sep by line
func text_plain(v interface{}) ([]byte, error) {
	conv_map := func(v map[string]interface{}) string {
		s := make([]string, 0, len(v))
		for k, v := range v {
			s = append(s, fmt.Sprintf("%v=%v", k, v))
		}
		return strings.Join(s, "\n")
	}

	conv_keyval := func(v map[string]string) string {
		s := make([]string, 0, len(v))
		for k, v := range v {
			s = append(s, fmt.Sprintf("%v=%v", k, v))
		}
		return strings.Join(s, "\n")
	}

	switch v := v.(type) {
	case map[string]interface{}:
		return []byte(conv_map(v)), nil
	case map[string]string:
		return []byte(conv_keyval(v)), nil
	default:
		return []byte(fmt.Sprintf("%v", v)), nil
	}
}
