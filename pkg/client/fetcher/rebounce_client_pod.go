package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
)

func (f *Fetcher) RebounceClientPod(version service.Version, serviceId string) {
	t := time.Now()
	log.Debugf("SudoryClientPod Rebounce: start")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	up := service.CreateUpdateService(version, serviceId, 1, 0, service.StepStatusProcessing, "", t, time.Time{})
	if err := f.sudoryAPI.UpdateServices(ctx, service.ConvertServiceStepUpdateClientToServer(up)); err != nil {
		log.Errorf("SudoryClientPod Rebounce: failed to update service status(processing): error: %s\n", err.Error())
	}

	f.ticker.Stop()
	log.Debugf("SudoryClientPod Rebounce: polling stop")

	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		log.Debugf("SudoryClientPod Rebounce: waiting to process the remaining services(timeout:30s)")

		for {
			<-time.After(time.Second * 3)
			services := f.RemainServices()
			if len(services) == 0 {
				break
			}

			buf := bytes.Buffer{}
			buf.WriteString("SudoryClientPod Rebounce: waiting remain services:")
			for uuid, status := range services {
				buf.WriteString(fmt.Sprintf("\n\tuuid: %s, status: %s", uuid, status.String()))
			}
			log.Debugf(buf.String() + "\n")
		}
	}()

	select {
	case <-time.After(time.Second * 30):
		log.Debugf("SudoryClientPod Rebounce: waiting timeout")
	case <-done:
		log.Debugf("SudoryClientPod Rebounce: waiting done")
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel2()

	up = service.CreateUpdateService(version, serviceId, 1, 0, service.StepStatusSuccess, "SudoryClient pod rebounce will be complete", t, time.Now())
	if err := f.sudoryAPI.UpdateServices(ctx2, service.ConvertServiceStepUpdateClientToServer(up)); err != nil {
		log.Errorf("SudoryClientPod Rebounce: failed to update service status(success): error: %s\n", err.Error())
	}

	f.Cancel()
}
