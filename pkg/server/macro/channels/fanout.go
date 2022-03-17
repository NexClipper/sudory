package channels

import "sync"

type Fanout struct {
	init sync.Once

	pub  *SafeChannel
	subs map[*SafeChannel]struct{}

	sync.Once
	sync.Mutex

	stop    func()
	closing chan struct{}
	closed  chan struct{}
}

func NewFanout(pub *SafeChannel) *Fanout {
	return &Fanout{pub: pub}
}

func (fanout *Fanout) SubChan() <-chan *SafeChannel {
	//lock gaurd
	fanout.Lock()
	defer fanout.Unlock()

	if false {
		//un
		for iter := range fanout.subs {
			if iter.IsClosed() {
				delete(fanout.subs, iter)
			}
		}
	}

	c := make(chan *SafeChannel, len(fanout.subs))
	defer close(c)

	for iter := range fanout.subs {
		c <- iter
	}

	return c
}

func (fanout *Fanout) NewSubscriber(size ...int) *SafeChannel {
	//lock gaurd
	fanout.Lock()
	defer fanout.Unlock()

	sub := NewSafeChannel(size...)
	Subscribe(fanout, sub)
	return sub
}

func (fanout *Fanout) Subscribe(sub ...*SafeChannel) {
	//lock gaurd
	fanout.Lock()
	defer fanout.Unlock()

	Subscribe(fanout, sub...)
}

func (fanout *Fanout) Unsubscribe(sub *SafeChannel) {
	//lock gaurd
	fanout.Lock()
	defer fanout.Unlock()
	Unsubscribe(fanout, sub)
}

func (fanout *Fanout) Stop() {
	if fanout.stop == nil {
		return
	}
	fanout.stop()
}

func (fanout *Fanout) Start() {
	fanout.Once.Do(func() {
		fanout.closing = make(chan struct{})
		fanout.closed = make(chan struct{})

		fanout.stop = func() {
			select {
			case fanout.closing <- struct{}{}:
				<-fanout.closed
			case <-fanout.closed:
			}

			fanout.Once = sync.Once{}
		}

		go func() {
			defer func() {
				fanout.closed <- struct{}{}

				close(fanout.closing)
				close(fanout.closed)
			}()

			for {
				select {
				case <-fanout.closing:
					return
				default:
				}

				select {
				case <-fanout.closing:
					return
				case v := <-fanout.pub.C:
					for sub := range fanout.SubChan() {
						sub.SafeSend(v) //fanout
					}
				}
			}
		}()
	})
}

func Subscribe(fanout *Fanout, sub ...*SafeChannel) {
	fanout.init.Do(func() {
		if fanout.subs == nil {
			fanout.subs = map[*SafeChannel]struct{}{}
		}
	})
	for _, sub := range sub {
		fanout.subs[sub] = struct{}{}
	}
}

func Unsubscribe(fanout *Fanout, sub ...*SafeChannel) {
	fanout.init.Do(func() {
		if fanout.subs == nil {
			fanout.subs = map[*SafeChannel]struct{}{}
		}
	})
	for _, sub := range sub {
		delete(fanout.subs, sub)
	}
}
