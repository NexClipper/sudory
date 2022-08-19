package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
)

func (f *Fetcher) Shutdown(serviceId string) {
	t := time.Now()
	log.Debugf("sudoryclient.shutdown: start")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := f.sudoryAPI.UpdateServices(ctx, &servicev2.HttpReq_ClientServiceUpdate{
		Uuid:     serviceId,
		Sequence: 0,
		Status:   servicev2.StepStatusProcessing,
		Started:  t,
	}); err != nil {
		log.Errorf("sudoryclient.shutdown: failed to update service status(processing): error: %s\n", err.Error())
	}

	f.ticker.Stop()
	log.Debugf("sudoryclient.shutdown: polling stop")

	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		log.Debugf("sudoryclient.shutdown: waiting to process the remaining services(timeout:30s)")

		for {
			<-time.After(time.Second * 3)
			services := f.RemainServices()
			if len(services) == 0 {
				break
			}

			buf := bytes.Buffer{}
			buf.WriteString("sudoryclient.shutdown: remain services:")
			for uuid, status := range services {
				buf.WriteString(fmt.Sprintf("\n\tuuid: %s, status: %s", uuid, status.String()))
			}
			log.Debugf(buf.String() + "\n")
		}
	}()

	select {
	case <-time.After(time.Second * 30):
		log.Debugf("sudoryclient.shutdown: timeout")
	case <-done:
		log.Debugf("sudoryclient.shutdown: done")
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel2()

	if err := f.sudoryAPI.UpdateServices(ctx2, &servicev2.HttpReq_ClientServiceUpdate{
		Uuid:     serviceId,
		Sequence: 0,
		Status:   servicev2.StepStatusSuccess,
		Result:   "sudoryclient shutdown will be complete",
		Started:  t,
		Ended:    time.Now(),
	}); err != nil {
		log.Errorf("sudoryclient.shutdown: failed to update service status(success): error: %s\n", err.Error())
	}

	f.Cancel()
}
