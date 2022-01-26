package channels

import (
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

// Distribute
//  sender(1) -> reciver(n,m)
//  pub: input channel
//  notify: call notify when the data is received
//  n: reciver gorutine count
//  m: reciver channel buffer size
//  stop: stop
func Distribute(pub *SafeChannel, notify func(v interface{}), n ...int) (stop func()) {

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

	closing := make(chan struct{})
	closed := make(chan struct{})

	wg := sync.WaitGroup{}

	//set stop
	stop = func() {
		select {
		case closing <- struct{}{}:
			<-closed
		case <-closed:
		}
		wg.Wait() //wg wait
	}

	go func() {

		//new reciver
		reciver := NewSafeChannel(chan_size)

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
		wg.Add(gorutine_cnt) // wg add
		for n := 0; n < gorutine_cnt; n++ {
			go func() {
				defer func() {
					wg.Done() //wg done
				}()
				for v := range reciver.C {
					notify(v) //notify
				}
			}()
		}
	}()
	return
}
