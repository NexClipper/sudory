package channels_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/channels"
)

func TestNewFanout(t *testing.T) {

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		<-time.After(time.Second * 3)
	}()

	pub := channels.NewSafeChannel(0)

	fanout := channels.NewFanout(pub)
	sub1 := channels.NewSafeChannel(0)
	fanout.Subscribe(sub1)
	sub2 := fanout.NewSubscriber()
	go func() {
		defer fanout.Unsubscribe(sub2)
		<-time.After(time.Second * 1)
	}()
	fanout.Start() //start

	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case <-ctx.Done():
				return
			case v := <-sub1.C:
				t.Log("sub1", v)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case <-ctx.Done():
				return
			case v := <-sub2.C:
				t.Log("sub2", v)
			}
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			defer func() {
				<-time.After(time.Millisecond * 100)
			}()
			pub.SafeSend(i)

			if i%2 == 0 {
				fanout.Stop()
				fanout.Start()
			}

		}
	}()

	wg.Wait()

	t.Log("done")
}
