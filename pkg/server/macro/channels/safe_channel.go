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

func WrapSafeChannel(c chan T) *SafeChannel {
	return &SafeChannel{C: c}
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
