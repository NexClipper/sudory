package gate

import (
	"sync"
	"sync/atomic"
)

// // closedchan is a reusable closed channel.
// var closedchan = channels.NewSafeChannel(0)

// func init() {
// 	closedchan.SafeClose()
// }

type Waiter interface {
	Wait() <-chan struct{}
}

type ManualReset struct {
	sync.Mutex

	// islocked bool
	channel atomic.Value
}

// NewManualReset
//  set: init status
func NewManualReset(set bool) *ManualReset {
	gate := &ManualReset{}
	// gate.islocked = set
	d := NewSafeChannel()
	if set {
		d.SafeClose()
	}

	gate.channel.Store(d)

	return gate
}

func (gate *ManualReset) Set() {
	//lock gaurd
	gate.Lock()
	defer gate.Unlock()

	d, _ := gate.channel.Load().(*SafeChannel)
	if !d.IsClosed() {
		d.SafeClose()
	}
	gate.channel.Store(d)

	// gate.islocked = false
}

func (signal *ManualReset) Reset() {
	//lock gaurd
	signal.Lock()
	defer signal.Unlock()

	d, _ := signal.channel.Load().(*SafeChannel)
	if d.IsClosed() {
		d = NewSafeChannel()
	}
	signal.channel.Store(d)

	// gate.islocked = true
}

func (signal *ManualReset) Wait() <-chan Empty {
	//lock gaurd
	signal.Lock()
	defer signal.Unlock()

	d, _ := signal.channel.Load().(*SafeChannel)

	return d.C
}

type AutoReset struct {
	sync.Mutex

	// islocked bool
	channel atomic.Value
}

func NewAutoReset(set bool) *AutoReset {
	gate := &AutoReset{}

	d := NewSafeChannel()
	if set {
		d.SafeClose()
	}

	gate.channel.Store(d)

	return gate
}

func (gate *AutoReset) Set() {
	//lock gaurd
	gate.Lock()
	defer gate.Unlock()

	d, _ := gate.channel.Load().(*SafeChannel)
	if !d.IsClosed() {
		d.SafeClose()
	}
	gate.channel.Store(d)

	// gate.islocked = false
}

func (gate *AutoReset) reset() {
	//lock gaurd
	gate.Lock()
	defer gate.Unlock()

	d, _ := gate.channel.Load().(*SafeChannel)
	if d.IsClosed() {
		d = NewSafeChannel()
	}
	gate.channel.Store(d)

	// gate.islocked = true
}

func (gate *AutoReset) Wait() <-chan Empty {
	//lock gaurd
	gate.Lock()
	defer gate.Unlock()

	d, _ := gate.channel.Load().(*SafeChannel)
	c := make(chan Empty)
	go func() {
		defer func() {
			gate.reset()
		}()
		<-d.C
		close(c)
	}()

	return c
}
