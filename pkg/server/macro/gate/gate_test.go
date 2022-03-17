package gate

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewManualReset(t *testing.T) {

	gate := NewManualReset(true)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	var seq int

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-gate.Wait():
				gate.Reset()
			}

			t.Log("manual reset gate test", seq)

			<-time.After(time.Millisecond * 1)
		}
	}()

	for i := 0; i < 10; i++ {

		seq = i
		gate.Set()

		<-time.After(time.Millisecond * 3)

	}
	cancel()

	wg.Wait()

	t.Log("done")
}

func TestNewAutoReset(t *testing.T) {

	gate := NewAutoReset(true)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	var seq int

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-gate.Wait():
			}

			t.Log("auto reset gate test", seq)

			<-time.After(time.Millisecond * 1)
		}
	}()

	for i := 0; i < 10; i++ {

		seq = i
		gate.Set()

		<-time.After(time.Millisecond * 3)

	}
	cancel()

	wg.Wait()

	t.Log("done")
}
