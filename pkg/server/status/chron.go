package status

import (
	"fmt"
	"io"
	"time"
)

type ChronUpdater interface {
	Update() error
}

func NewChron(ErrOut io.Writer, interval time.Duration, chron ...func() ChronUpdater) func() {

	new_chrons := make([]ChronUpdater, len(chron))
	for n := range chron {
		new_chrons[n] = chron[n]()
	}

	tick := time.NewTicker(interval)

	closing := make(chan interface{})
	closed := make(chan interface{})

	//set stop
	stop := func() {
		select {
		case closing <- nil:
			<-closed
		case <-closed:
		}
	}

	go func() {

		defer func() {
			close(closed)
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
			case <-tick.C:
				for n := range new_chrons {
					if err := new_chrons[n].Update(); err != nil {
						fmt.Fprintf(ErrOut, "%s", fmt.Errorf("chron error: %w", err).Error())
					}
				}
			}
		}
	}()

	return func() {
		tick.Stop()
		stop()
	}
}
