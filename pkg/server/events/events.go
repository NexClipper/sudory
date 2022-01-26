package events

// EventArgs
//  이벤트 전달 인자
type EventArgs struct {
	Sender string
	Args   map[string]interface{}
}

// Invoker
type Invoker func(*EventArgs)

// Invoke
//  이벤트 입력
//  활성 하지 않으면 초기값은 default invoker
var Invoke Invoker = default_invoker

//default invoker (do nothing)
func default_invoker(*EventArgs) {
	println("this event message recived in default invoker :)")
}

//Activate 함수에서 사용
// var once sync.Once

//이벤트 매니저
// var Manager *EventManager

// // 이벤트 매니저 활성
// func Activate(sender *channels.SafeChannel, events []EventContexter) (*EventManager, func()) {

// 	// once.Do(func() {
// 	mngr := NewManager(log.Printf)
// 	defer func() {
// 		mngr.Events = events //setting EventContexter
// 		// Manager = mngr       //setting global manager
// 	}()
// 	// sender := channels.NewSafeChannel(0) //sender
// 	notify := func(v interface{}) {
// 		args, _ := v.(*EventArgs) //입력 받은 이벤트 데이터, 캐스트 해본다
// 		mngr.Notify(args)
// 	}
// 	reciver_count := runtime.NumCPU() //reciver count

// 	stop := channels.Distribute(sender, notify, reciver_count) //setting Distribute
// 	stop = func() {
// 		//채널 종료
// 		stop()

// 		//이벤트 종료 기다림
// 		for _, event := range mngr.Events {
// 			event_, ok := event.(Waiter)
// 			if ok {
// 				event_.Wait()
// 			}
// 		}

// 		//파일 리스너를 위한, 파일 핸들러 종료
// 		Files.CloseFileAll()
// 	}

// 	// })
// 	return mngr, stop
// }
