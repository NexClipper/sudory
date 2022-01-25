package poll

import (
	"encoding/json"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/panta/machineid"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

const defaultPollingInterval = 5 // * time.Second

type Poller struct {
	token            string
	server           string
	machineID        string
	clusterId        string
	client           *httpclient.HttpClient
	pollingScheduler *gocron.Scheduler
	serviceScheduler *service.ServiceScheduler
}

func NewPoller(token, server, clusterId string, serviceScheduler *service.ServiceScheduler) (*Poller, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}

	return &Poller{token: token, server: server, machineID: id, clusterId: clusterId, client: httpclient.NewHttpClient(server, token, 0, 0), pollingScheduler: gocron.NewScheduler(time.UTC), serviceScheduler: serviceScheduler}, nil
}

func (p *Poller) Start() {
	p.pollingScheduler.Every(defaultPollingInterval).Second().Do(p.poll)
	p.pollingScheduler.StartAsync()
}

func (p *Poller) ChangePollingInterval(interval int) {
	p.pollingScheduler.Clear()
	p.pollingScheduler.Every(interval).Second().Do(p.poll)
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

	body, err := p.client.PutJson("/client/service", map[string]string{"cluster_uuid": p.clusterId}, jsonb)
	if err != nil {
		p.serviceScheduler.RepairUpdateFailedServices(updatedServices)
		log.Errorf(err.Error())
		return
	}

	respData := []servicev1.HttpRspClientSideService{}
	if err := json.Unmarshal(body, &respData); err != nil {
		log.Errorf(err.Error())
		return
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
