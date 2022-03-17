package gate

import "sync"

type Empty struct{}

type SafeChannel struct {
	C      chan Empty
	closed bool
	mux    sync.Mutex
}

func NewSafeChannel() *SafeChannel {
	return &SafeChannel{C: make(chan Empty)}
}

func (channel *SafeChannel) IsClosed() bool {
	channel.mux.Lock()
	defer channel.mux.Unlock()
	return channel.closed
}

func (channel *SafeChannel) SafeClose() {
	channel.mux.Lock()
	defer channel.mux.Unlock()
	if !channel.closed {
		close(channel.C)
		channel.closed = true
	}
}

func (channel *SafeChannel) Set() bool {
	channel.mux.Lock()
	defer channel.mux.Unlock()
	if !channel.closed {
		channel.C <- Empty{}
		return true
	}
	return false
}
