package events

import (
	"runtime"
	"sync"
)

type T interface{}

type SafeChannel struct {
	C      chan T
	closed bool
	mux    sync.Mutex
}

func NewSafeChannel(size ...int) *SafeChannel {
	var buff_size int = 0
	if 0 < len(size) {
		buff_size = size[0]
	}
	return &SafeChannel{C: make(chan T, buff_size)}
}

func (me *SafeChannel) IsClosed() bool {
	me.mux.Lock()
	defer me.mux.Unlock()
	return me.closed
}

func (me *SafeChannel) SafeClose() {
	me.mux.Lock()
	defer me.mux.Unlock()
	if !me.closed {
		close(me.C)
		me.closed = true
	}
}

func (me *SafeChannel) SafeSend(value T) bool {
	me.mux.Lock()
	defer me.mux.Unlock()
	if !me.closed {
		me.C <- value
		return true
	}
	return false
}

func PubSub(pub *SafeChannel, notify func(v interface{}), size ...int) (stop func()) {

	closing := make(chan struct{})
	closed := make(chan struct{})

	//set stop
	stop = func() {
		select {
		case closing <- struct{}{}:
			<-closed
		case <-closed:
		}
	}

	go func() {

		//new reciver
		reciver := NewSafeChannel(size...)

		//sender
		go func() {
			defer func() {
				close(closed)
				reciver.SafeClose()
			}()

			for {
				select {
				case <-closing:
					return
				default:
				}

				select {
				case <-closing:
					return
				case v := <-pub.C:
					reciver.SafeSend(v)
				}
			}
		}()

		//reciver
		num_cpu := runtime.NumCPU()
		wg := sync.WaitGroup{}
		wg.Add(num_cpu) // wg add
		for n := 0; n < num_cpu; n++ {
			go func() {
				defer wg.Done() //wg done

				for v := range reciver.C {
					if args, ok := v.(map[string]interface{}); ok {
						notify(args) //notify
					}
				}
			}()
		}

		wg.Wait() //wg wait
	}()
	return
}
