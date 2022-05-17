package event

import (
	"context"
	"sync"
	"time"

	"github.com/NexClipper/sudory/pkg/server/event/fifo"
	"github.com/NexClipper/sudory/pkg/server/macro/gate"
	"github.com/pkg/errors"
)

type EventNotifierMux struct {
	config EventNotifierMuxerConfig
	pub    EventPublisher

	notifiers    HashsetNotifier
	errorHandler HashsetErrorHandlers

	*fifo.Queue //arrival item queue

	onceforupdater sync.Once         //once for notifiy routine
	closing        context.Context   //subscriber closing context
	doClose        func()            //subscriber closing caller
	closed         *gate.ManualReset //subscriber closed gate
	arrival        *gate.ManualReset //item arrival gate
}

var _ EventNotifiMuxConfigHolder = (*EventNotifierMux)(nil)
var _ EventNotifierMuxer = (*EventNotifierMux)(nil)

func NewEventSubscribe(cfg EventNotifierMuxerConfig, errorHandler HashsetErrorHandlers) *EventNotifierMux {
	if cfg.UpdateInterval == time.Duration(0) {
		cfg.UpdateInterval = time.Second * 15 //default update-interval is 15s
	}

	sub := &EventNotifierMux{}
	sub.notifiers = HashsetNotifier{}
	sub.Queue = fifo.NewQueue()
	sub.closing, sub.doClose = context.WithCancel(context.Background())
	sub.closed = gate.NewManualReset(true)
	sub.arrival = gate.NewManualReset(false)

	sub.config = cfg //config
	sub.errorHandler = errorHandler

	return sub
}

func (mux EventNotifierMux) Config() *EventNotifierMuxerConfig {
	return &mux.config
}

func (mux EventNotifierMux) Notifiers() HashsetNotifier {
	return mux.notifiers
}

func (mux *EventNotifierMux) Update(sender string, v ...interface{}) {
	//이벤트 이름으로 처리하는 이벤트 검사
	if mux.config.Name != sender {
		return //처리하는 이벤트가 아니다
	}

	for _, v := range v {
		if v == nil {
			continue
		}

		//add arrivals in queue
		mux.Queue.Add(v)
	}

	//release arrival wait
	mux.arrival.Set()

	mux.onceforupdater.Do(func() {
		//goroutine 생성하면서, closed 게이트 활성
		//생성 함수에서 미리 활성 하면,
		//Update 함수를 타지 않는 경우 closed 게이트를 풀어 주는 부분이 없어서
		//closed를 기다리는 mux.Close() 함수에서 무한 대기
		mux.closed.Reset()

		go func() {
			defer func() {
				//goroutine 종료하면서, closed 게이트 풀어준다
				//이어서 goroutine 종료한 mux.Close() 함수에서
				//대기를 완료하고 subscriber가 종료된다.
				mux.closed.Set()
			}()
		LOOP:
			for {
				select {
				case <-mux.closing.Done(): //wait closing
					//큐가 비어있어야 끝난다
					if mux.Queue.Len() == 0 {
						return
					}
				case <-mux.arrival.Wait(): //wait arrival
					mux.arrival.Reset()
				case <-time.After(mux.config.UpdateInterval): //wait timer
				}

				//check queue length
				queue_length := mux.Queue.Len()
				if queue_length == 0 {
					continue LOOP //queue is empty
				}

				sl := make([]interface{}, 0, queue_length)
				for i := 0; i < queue_length; i++ {
					item := mux.Queue.Next()
					sl = append(sl, item)
				}
				factory := NewMarshalFactory(sl...)

				//모든 리스너의 Update 호출 (async)
				futures := mux.notifiers.OnNotifyAsync(factory)

				for _, future := range futures {
					for future := range future {
						if future.Error != nil {
							//업데이트 오류 처리
							mux.errorHandler.OnError(errors.Wrapf(future.Error, "event notify %s",
								MapString(future.Notifier.Property())))
						}
					}
				}

			} //! LOOP:
		}() //! go func()
	})
}

func (mux *EventNotifierMux) Regist(pub EventPublisher) {
	if mux.pub != nil {
		mux.pub = pub
	}
	if mux.pub != nil {
		mux.pub.NotifierMuxers().Add(mux) //Subscribe
	}
}

func (mux *EventNotifierMux) Close() {
	//unsubscribe
	if mux.pub != nil {
		mux.pub.NotifierMuxers().Remove(mux) //Unsubscribe
	}

	//close the goroutine
	if mux.doClose != nil {
		mux.doClose()
	}
	//wait for terminate the goroutine
	//
	<-mux.closed.Wait()

	//clear notifiers
	for iter := range mux.Notifiers() {
		iter.Close()
	}
}
