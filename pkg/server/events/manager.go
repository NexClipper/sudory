package events

import (
	"log"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/channels"
)

// ErrorHandler
type ErrorHandler func(fmt string, v ...interface{})

// EventManager
type EventManager struct {
	// Stop         context.CancelFunc

	errorHandler  ErrorHandler
	eventContexts []EventContexter
	sender        *channels.SafeChannel
	Invoker       func(v *EventArgs)
}

// NewManager
func NewManager(sender *channels.SafeChannel, eventContexts []EventContexter, handler ErrorHandler) *EventManager {
	if handler == nil {
		handler = log.Printf //default logger
	}
	if sender == nil {
		log.Fatal("invalid parameter sender")
		return nil
	}

	invoke := func(v *EventArgs) { sender.SafeSend(v) }

	return &EventManager{
		errorHandler:  handler,
		eventContexts: eventContexts,
		sender:        sender,
		Invoker:       invoke,
	}
}

// notify
//  이밴트 리스너로 전달
func (manager EventManager) notify(args *EventArgs) {
	for _, ectx := range manager.eventContexts {

		//값 가공
		args_ := map[string]interface{}{
			"name":  ectx.Name(), //이벤트 이름 추가
			"issue": time.Now(),  //시간 추가
			"args":  args.Args,
		}

		ectx.Raise(EventArgs{Sender: args.Sender, Args: args_}, manager.errorHandler)
	}

}

func Activate(manager *EventManager, n ...int) func() {

	// reciver gorutine count
	//  0보다 커야함
	var gorutine_cnt int = 1
	if 0 < len(n) {
		gorutine_cnt = n[0]
	}
	if gorutine_cnt <= 0 {
		gorutine_cnt = 1 //0보다 커야함
	}

	// reciver channel buffer size
	//  음수는 안됨
	var chan_size int = 0
	if 1 < len(n) {
		chan_size = n[1]
	}
	if chan_size < 0 {
		chan_size = 0 //음수면 안됨
	}

	notify := func(v interface{}) {
		args, _ := v.(*EventArgs) //입력 받은 이벤트 데이터, 캐스트
		manager.notify(args)
	}

	stop := channels.Distribute(manager.sender, notify, gorutine_cnt, chan_size) //setting Distribute

	closed := channels.NewSafeChannel(0)  //closed signal
	closing := channels.NewSafeChannel(0) //closing signal
	go func() {
		defer func() { closing.SafeClose() }()
		<-closing.C //wait closing signal

		//채널 종료
		stop()

		//이벤트 종료 기다림
		for _, event := range manager.eventContexts {
			event_, ok := event.(Waiter)
			if ok {
				event_.Wait()
			}
		}

		//파일 리스너를 위한, 파일 핸들러 종료
		Files.CloseFileAll()

		closed.SafeSend(nil) //send closed signal
	}()

	return func() {
		defer func() { closed.SafeClose() }()
		closing.SafeSend(nil) //send closing signal
		<-closed.C            //wait closed signal
	}
}
