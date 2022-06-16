package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/panta/machineid"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/scheduler"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

const (
	defaultPollingInterval  = 5 // * time.Second
	minPollingInterval      = 5
	maxRetryForFailedToSend = 5
)

var (
	failedToSendCountMap sync.Map
)

type Fetcher struct {
	bearerToken     string
	machineID       string
	clusterId       string
	client          *httpclient.HttpClient
	ticker          *time.Ticker
	pollingInterval int
	scheduler       *scheduler.Scheduler
	done            chan struct{}
}

func NewFetcher(bearerToken, server, clusterId string, sch *scheduler.Scheduler) (*Fetcher, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	id = strings.ReplaceAll(id, "-", "")

	return &Fetcher{
		bearerToken:     bearerToken,
		machineID:       id,
		clusterId:       clusterId,
		client:          httpclient.NewHttpClient(server, "", 0, 0),
		ticker:          time.NewTicker(defaultPollingInterval * time.Second),
		pollingInterval: defaultPollingInterval,
		scheduler:       sch,
		done:            make(chan struct{})}, nil
}

func (f *Fetcher) ChangePollingInterval(interval int) (int, error) {
	if interval < minPollingInterval {
		return 0, fmt.Errorf("interval(%d) you want to change is less than the minimum interval(%d)", interval, minPollingInterval)
	}

	if f.pollingInterval == interval {
		return 0, nil
	}

	f.pollingInterval = interval
	f.ticker.Reset(time.Second * time.Duration(interval))

	return interval, nil
}

func (f *Fetcher) Done() <-chan struct{} {
	return f.done
}

func (f *Fetcher) Cancel() {
	close(f.done)
}

func (f *Fetcher) Polling(ctx context.Context) error {
	if f == nil || f.ticker == nil {
		return fmt.Errorf("fetcher or fetcher.ticker is not created")
	}

	go f.UpdateServiceProcess()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-f.ticker.C:
				f.poll()
			}
		}
	}()

	return nil
}

func (f *Fetcher) poll() {
	// check if there are any services that have failed to send
	failedServices := f.scheduler.GetServicesWithFailedToSendFlag()

	for _, serv := range failedServices {
		failedCnt := 0
		if cnt, ok := failedToSendCountMap.Load(serv.GetServiceId()); ok {
			failedCnt = cnt.(int)
		}

		if failedCnt > maxRetryForFailedToSend {
			failedToSendCountMap.Delete(serv.GetServiceId())
			f.scheduler.DeleteServiceWithFlag(serv.GetServiceId(), scheduler.ServiceCheckedFlagFailedToSend)
			continue
		}

		// services -> reqData
		reqData := scheduler.ServiceClientToServer(serv)

		if reqData == nil {
			continue
		}

		jsonb, err := json.Marshal(reqData)
		if err != nil {
			log.Errorf(err.Error())
			failedToSendCountMap.Store(serv.GetServiceId(), failedCnt+1)
			continue
		}

		if _, err := f.client.PutJson("/client/service", nil, jsonb); err != nil {
			log.Errorf(err.Error())
			failedToSendCountMap.Store(serv.GetServiceId(), failedCnt+1)
			continue
		}
		f.scheduler.DeleteServiceWithFlag(serv.GetServiceId(), scheduler.ServiceCheckedFlagFailedToSend)
	}

	body, err := f.client.GetJson("/client/service", nil)
	if err != nil {
		log.Errorf(err.Error())

		if f.client.IsTokenExpired() {
			f.ticker.Stop()
			go f.RetryHandshake()
		}

		return
	}

	f.ChangeClientConfigFromToken()

	respData := []servicev1.HttpRspService_ClientSide{}
	if body != nil {
		if err := json.Unmarshal(body, &respData); err != nil {
			log.Errorf(err.Error())
			return
		}
	}
	log.Debugf("Recived %d service from server.", len(respData))

	if len(respData) == 0 {
		return
	}

	// respData -> services
	recvServices := scheduler.ServiceListServerToClient(respData)

	// Register new services.
	f.scheduler.RegisterServices(recvServices)
}

func (f *Fetcher) ChangeClientConfigFromToken() {
	claims := new(sessionv1.ClientSessionPayload)
	jwt_token, _, err := jwt.NewParser().ParseUnverified(f.client.GetToken(), claims)
	if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || err != nil {
		log.Warnf("Failed to bind payload : %v\n", err)
		return
	}

	if interval, err := f.ChangePollingInterval(claims.PollInterval); err != nil {
		log.Warnf("Failed to change polling interval : %v\n", err)
	} else {
		if interval != 0 {
			log.Debugf("Change polling interval to %v\n", claims.PollInterval)
		}
	}

	if log.GetLogger().GetLevel() != strings.ToLower(claims.Loglevel) {
		if err := log.GetLogger().SetLevel(claims.Loglevel); err == nil {
			log.Debugf("Changed logger level to %s\n", claims.Loglevel)
		}
	}
}

func (f *Fetcher) UpdateServiceProcess() {
	for serv := range f.scheduler.NotifyServiceUpdate() {
		if serv == nil {
			continue
		}

		jsonb, err := json.Marshal(serv)
		if err != nil {
			log.Errorf(err.Error())
			f.scheduler.ChangeServiceFlagFailedToSend(serv.Uuid)
			continue
		}

		if _, err := f.client.PutJson("/client/service", nil, jsonb); err != nil {
			log.Errorf(err.Error())
			f.scheduler.ChangeServiceFlagFailedToSend(serv.Uuid)
			continue
		}
		f.scheduler.DeleteServiceWithFlag(serv.Uuid, scheduler.ServiceCheckedFlagDone)
	}
}
