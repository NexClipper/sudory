package ticker

import (
	"reflect"
	"time"
)

type HashsetErrorHandlers map[uintptr]func(error)

func (hashset HashsetErrorHandlers) Add(fn ...func(error)) {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		hashset[ptr] = fn
	}
}
func (hashset HashsetErrorHandlers) Remove(fn ...func(error)) {
	for _, fn := range fn {
		ptr := reflect.ValueOf(fn).Pointer()
		delete(hashset, ptr)
	}
}
func (hashset HashsetErrorHandlers) OnError(err error) {
	for _, handler := range hashset {
		handler(err)
	}
}

func NewTicker(interval time.Duration, fn ...func()) func() {
	tick := time.NewTicker(interval)

	//set stop
	closing := make(chan interface{})
	closed := make(chan interface{})
	stop := func() {
		select {
		case closing <- nil:
			<-closed
		case <-closed:
		}
	}
	close_ := func() {
		tick.Stop()
		stop()
	}

	go func() {
		defer func() {
			close(closed)
		}()

		for {
			select {
			case <-closing:
				return
			case <-tick.C:
			}

			for _, fn := range fn {
				fn()
			}
		}
	}()

	return close_
}
