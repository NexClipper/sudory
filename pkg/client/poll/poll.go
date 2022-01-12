package poll

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/panta/machineid"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/service"
	"github.com/NexClipper/sudory/pkg/server/model"
)

const defaultPollingInterval = 5 // * time.Second

type Poller struct {
	token            string
	server           string
	machineID        string
	client           *httpclient.HttpClient
	pollingScheduler *gocron.Scheduler
	serviceScheduler *service.ServiceScheduler
}

func NewPoller(token, server string, serviceScheduler *service.ServiceScheduler) *Poller {
	id, err := machineid.ID()
	if err != nil {
		return nil
	}

	//log.Printf("machine id: %s", id)

	uri := server + "/client/service"

	return &Poller{token: token, server: server, machineID: id, client: httpclient.NewHttpClient(uri, token), pollingScheduler: gocron.NewScheduler(time.UTC), serviceScheduler: serviceScheduler}
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
	// Get services's status
	// servicesWillUpdate := p.serviceScheduler.GetServices()

	// TODO: services' status -> reqData

	reqData := &model.ReqClientGetService{
		ClusterID: 999,
	}

	body, err := p.client.PutJson(reqData)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// TODO: If server updated service's status, remove completed services.

	respData := &model.RespService{}
	if err := json.Unmarshal(body, respData); err != nil {
		log.Printf(err.Error())
		return
	}

	// TODO: respData -> services
	recvServices := make(map[string]*service.Service)

	// Register new services.
	p.serviceScheduler.RegisterServices(recvServices)
}
