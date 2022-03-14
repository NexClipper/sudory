package fetcher

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/panta/machineid"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/scheduler"
	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

const (
	defaultPollingInterval = 5 // * time.Second
	minPollingInterval     = 5
)

type Fetcher struct {
	bearerToken     string
	machineID       string
	clusterId       string
	client          *httpclient.HttpClient
	ticker          *time.Ticker
	pollingInterval int
	scheduler       *scheduler.Scheduler
}

func NewFetcher(bearerToken, server, clusterId string, scheduler *scheduler.Scheduler) (*Fetcher, error) {
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
		scheduler:       scheduler}, nil
}

func (f *Fetcher) ChangePollingInterval(interval int) error {
	if f.pollingInterval == interval || interval < minPollingInterval {
		return fmt.Errorf("interval(%d) you want to change is the same as the previous interval(%d) or less than the minimum interval(%d)", interval, f.pollingInterval, minPollingInterval)
	}

	f.pollingInterval = interval
	f.ticker.Reset(time.Second * time.Duration(interval))

	return nil
}

func (f *Fetcher) Polling() error {
	if f == nil || f.ticker == nil {
		return fmt.Errorf("fetcher or fetcher.ticker is not created")
	}

	go func() {
		for ; ; <-f.ticker.C {
			f.poll()
		}
	}()

	return nil
}

func (f *Fetcher) poll() {
	// Get updated services. If the service is done, it is deleted.
	updatedServices := f.scheduler.GetServicesWithUpdatedDoneFlag()

	// services -> reqData
	reqData := scheduler.ServiceListClientToServer2(updatedServices)

	jsonb, err := json.Marshal(reqData)
	if err != nil {
		f.scheduler.RollbackServicesWithDoneUpdatedFlag(updatedServices)
		log.Errorf(err.Error())
	}

	body, err := f.client.PutJson("/client/service", nil, jsonb)
	if err != nil {
		f.scheduler.RollbackServicesWithDoneUpdatedFlag(updatedServices)
		log.Errorf(err.Error())
		return
	}
	f.scheduler.DeleteServicesWithDoneFlag(updatedServices)

	f.ChangeClientConfigFromToken()

	respData := []servicev1.HttpRspClientSideService{}
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
	recvServices := scheduler.ServiceListServerToClient2(respData)

	// Register new services.
	f.scheduler.RegisterServices(recvServices)
}

func (f *Fetcher) ChangeClientConfigFromToken() {
	cfgM := make(map[string]interface{})

	if err := jwt.BindPayload(f.client.GetToken(), &cfgM); err != nil {
		log.Warnf("Failed to bind payload : %v\n", err)
		return
	}

	if v, b := cfgM["poll-interval"]; b {
		if pollInterval, ok := v.(float64); !ok {
			log.Warnf("Failed to assert type for poll-interval(%v).", v)
		} else {
			if err := f.ChangePollingInterval(int(pollInterval)); err != nil {
				log.Warnf("Failed to change polling interval : %v\n", err)
			} else {
				log.Debugf("Change polling interval to %v\n", pollInterval)
			}
		}
	}

	if v, b := cfgM["Loglevel"]; b {
		if logLevel, ok := v.(string); !ok {
			log.Warnf("Failed to assert type for Loglevel(%v).\n", v)
		} else {
			if err := log.GetLogger().SetLevel(logLevel); err != nil {
				log.Warnf("Failed to change logger level : %v.\n", err)
			} else {
				log.Debugf("Changed logger level to %s\n", logLevel)
			}
		}
	}
}
