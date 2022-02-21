package poll

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/panta/machineid"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

const (
	defaultPollingInterval = 5 // * time.Second
	minPollingInterval     = 5
)

type Poller struct {
	bearerToken      string
	server           string
	machineID        string
	clusterId        string
	client           *httpclient.HttpClient
	pollingInterval  int
	pollingScheduler *gocron.Scheduler
	serviceScheduler *service.ServiceScheduler
}

func NewPoller(bearerToken, server, clusterId string, serviceScheduler *service.ServiceScheduler) (*Poller, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	id = strings.ReplaceAll(id, "-", "")

	return &Poller{bearerToken: bearerToken, server: server, machineID: id, clusterId: clusterId, client: httpclient.NewHttpClient(server, "", 0, 0), pollingInterval: defaultPollingInterval, pollingScheduler: gocron.NewScheduler(time.UTC), serviceScheduler: serviceScheduler}, nil
}

func (p *Poller) Start() {
	p.pollingScheduler.Every(p.pollingInterval).Second().Do(p.poll)
	p.pollingScheduler.StartAsync()
}

func (p *Poller) ChangePollingInterval(interval int) error {
	if p.pollingInterval == interval || interval < minPollingInterval {
		return fmt.Errorf("interval(%d) you want to change is the same as the previous interval(%d) or less than the minimum interval(%d)", interval, p.pollingInterval, minPollingInterval)
	}

	p.pollingInterval = interval
	p.pollingScheduler.Clear()
	p.pollingScheduler.Every(interval).Second().Do(p.poll)

	return nil
}

func (p *Poller) poll() {
	// Get updated services. If the service is done, it is deleted.
	updatedServices := p.serviceScheduler.GetDeleteServicesUpdated()

	// services -> reqData
	reqData := service.ServiceListClientToServer(updatedServices)

	jsonb, err := json.Marshal(reqData)
	if err != nil {
		p.serviceScheduler.RepairUpdateFailedServices(updatedServices)
		log.Errorf(err.Error())
	}

	body, err := p.client.PutJson("/client/service", nil, jsonb)
	if err != nil {
		p.serviceScheduler.RepairUpdateFailedServices(updatedServices)
		log.Errorf(err.Error())
		return
	}

	p.ChangeClientConfigFromToken()

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
	recvServices := service.ServiceListServerToClient(respData)
	// Delete duplicated services.
	recvServices = p.serviceScheduler.DeleteDuplicatedServices(recvServices)

	// Register new services.
	p.serviceScheduler.RegisterServices(recvServices)
}

func (p *Poller) ChangeClientConfigFromToken() {
	cfgM := make(map[string]interface{})

	if err := jwt.BindPayload(p.client.GetToken(), &cfgM); err != nil {
		log.Warnf("Failed to bind payload : %v\n", err)
		return
	}

	if v, b := cfgM["poll-interval"]; b {
		if pollInterval, ok := v.(float64); !ok {
			log.Warnf("Failed to assert type for poll-interval(%v).", v)
		} else {
			if err := p.ChangePollingInterval(int(pollInterval)); err != nil {
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
