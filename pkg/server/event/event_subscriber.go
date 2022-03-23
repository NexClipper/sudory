package event

import (
	"context"
	"sync"
	"time"

	"github.com/NexClipper/sudory/pkg/server/event/fifo"
	"github.com/NexClipper/sudory/pkg/server/macro/gate"
	"github.com/pkg/errors"
)

type EventSub struct {
	config EventSubscribeConfig
	pub    EventPublisher

	notifiers    HashsetNotifiers
	errorHandler HashsetErrorHandlers

	*fifo.Queue //arrival item queue

	onceforupdater sync.Once         //once for notifiy routine
	closing        context.Context   //subscriber closing context
	doClose        func()            //subscriber closing caller
	closed         *gate.ManualReset //subscriber closed gate
	arrival        *gate.ManualReset //item arrival gate
}

func NewEventSubscribe(cfg EventSubscribeConfig, errorHandler HashsetErrorHandlers) *EventSub {
	if cfg.UpdateInterval == time.Duration(0) {
		cfg.UpdateInterval = time.Second * 15 //default update-interval is 15s
	}

	sub := &EventSub{}
	sub.notifiers = HashsetNotifiers{}
	sub.Queue = fifo.NewQueue()
	sub.closing, sub.doClose = context.WithCancel(context.Background())
	sub.closed = gate.NewManualReset(true)
	sub.arrival = gate.NewManualReset(false)

	sub.config = cfg //config
	sub.errorHandler = errorHandler

	return sub
}

func (subscriber EventSub) Config() *EventSubscribeConfig {
	return &subscriber.config
}

func (subscriber EventSub) Notifiers() HashsetNotifiers {
	return subscriber.notifiers
}

// func (subscriber EventSub) ErrorHandlers() HashsetErrorHandlers {
// 	return subscriber.errorHandler
// }

func (subscriber *EventSub) Update(sender string, v ...interface{}) {
	//이벤트 이름으로 처리하는 이벤트 검사
	if subscriber.config.Name != sender {
		return //처리하는 이벤트가 아니다
	}

	for _, v := range v {
		if v == nil {
			continue
		}

		//add arrivals in queue
		subscriber.Queue.Add(v)
	}

	//release arrival wait
	subscriber.arrival.Set()

	subscriber.onceforupdater.Do(func() {
		//goroutine 생성하면서, closed 게이트 활성
		//생성 함수에서 미리 활성 하면,
		//Update 함수를 타지 않는 경우 closed 게이트를 풀어 주는 부분이 없어서
		//closed를 기다리는 subscriber.Close() 함수에서 무한 대기
		subscriber.closed.Reset()

		go func() {
			defer func() {
				//goroutine 종료하면서, closed 게이트 풀어준다
				//이어서 goroutine 종료한 subscriber.Close() 함수에서
				//대기를 완료하고 subscriber가 종료된다.
				subscriber.closed.Set()
			}()
		LOOP:
			for {
				select {
				case <-subscriber.closing.Done(): //wait closing
					//큐가 비어있어야 끝난다
					if subscriber.Queue.Len() == 0 {
						return
					}
				case <-subscriber.arrival.Wait(): //wait arrival
					subscriber.arrival.Reset()
				case <-time.After(subscriber.config.UpdateInterval): //wait timer
				}

				//check queue length
				queue_length := subscriber.Queue.Len()
				if queue_length == 0 {
					continue LOOP //queue is empty
				}

				sl := make([]interface{}, 0, queue_length)
				for i := 0; i < queue_length; i++ {
					item := subscriber.Queue.Next()
					sl = append(sl, item)
				}
				factory := NewMarshalFactory(sl...)

				//모든 리스너의 Update 호출 (async)
				futures := subscriber.notifiers.OnNotifyAsync(factory)

				for _, future := range futures {
					for future := range future {
						if future.Error != nil {
							//업데이트 오류 처리
							subscriber.errorHandler.OnError(errors.Wrapf(future.Error, "event notify %s",
								future.Notifier.PropertyString()))
						}
					}
				}

			} //! LOOP:
		}() //! go func()
	})
}

func (subscriber *EventSub) Regist(pub EventPublisher) {
	subscriber.pub = pub
	if subscriber.pub != nil {
		pub.Subscribers().Add(subscriber) //Subscribe
	}
}

func (subscriber *EventSub) Close() {
	//unsubscribe
	if subscriber.pub != nil {
		subscriber.pub.Subscribers().Remove(subscriber) //Unsubscribe
	}

	//close the goroutine
	subscriber.doClose()
	//wait for terminate the goroutine
	//
	<-subscriber.closed.Wait()

	//clear notifiers
	for iter := range subscriber.notifiers {
		iter.Close()
	}
}
