package scheduler

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/NexClipper/sudory/pkg/client/executor"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
)

const defaultMaxProcessLimit = 10

type Scheduler struct {
	servicesStatusMap map[string]service.ServiceStatus
	lock              sync.RWMutex
	maxProcessLimit   int
	updateChan        chan service.ServiceUpdateInterface // this channel receives service's status
	notifyUpdateChan  chan service.ServiceUpdateInterface
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		servicesStatusMap: make(map[string]service.ServiceStatus),
		maxProcessLimit:   defaultMaxProcessLimit,
		updateChan:        make(chan service.ServiceUpdateInterface),
		notifyUpdateChan:  make(chan service.ServiceUpdateInterface)}
}

func (s *Scheduler) Start() error {
	if s.updateChan == nil || s.servicesStatusMap == nil {
		return fmt.Errorf("scheduler don't have channel")
	}

	go s.RecvNotifyServiceStatus()

	return nil
}

func (s *Scheduler) RegisterServices(services map[string]service.ServiceInterface) {
	// 1. already existing services drop
	var startingList []service.ServiceInterface
	s.lock.Lock()
	for _, srv := range services {
		_, ok := s.servicesStatusMap[srv.GetId()]
		if !ok {
			startingList = append(startingList, srv)
		}
	}
	sort.Slice(startingList, func(i, j int) bool {
		if startingList[i].GetPriority() > startingList[j].GetPriority() {
			return true
		} else if startingList[i].GetPriority() < startingList[j].GetPriority() {
			return false
		} else {
			return startingList[i].GetCreatedTime().Before(startingList[j].GetCreatedTime())
		}
	})

	// 2. if existing service's status is ServiceStatusSuccess or ServiceStatusFailed, delete in statusMap
	var preExistingServiceUuids []string
	var deleteServiceUuids []string
	for uuid, status := range s.servicesStatusMap {
		preExistingServiceUuids = append(preExistingServiceUuids, uuid+" "+status.String())
		if status == service.ServiceStatusSuccess || status == service.ServiceStatusFailed {
			deleteServiceUuids = append(deleteServiceUuids, uuid)
			delete(s.servicesStatusMap, uuid)
		}
	}

	// 3. maxProcessLimit - len(statusMap) is number starting now
	remain := s.maxProcessLimit - len(s.servicesStatusMap)
	log.Debugf("Number of services that can be started : %d\n", remain)

	var startingServiceUuids []string
	for _, serv := range startingList {
		if remain > 0 {
			startingServiceUuids = append(startingServiceUuids, serv.GetId())
			s.servicesStatusMap[serv.GetId()] = service.ServiceStatusPreparing
			// create and execute(goroutine) service.
			go s.ExecuteService(serv)
			remain--
		} else {
			break
		}
	}

	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("\nPreExistingServices: %d", len(preExistingServiceUuids)))
	for _, uuid := range preExistingServiceUuids {
		buf.WriteString(fmt.Sprintf("\n\t%s", uuid))
	}

	buf.WriteString(fmt.Sprintf("\nDeleteServices: %d", len(deleteServiceUuids)))
	for _, uuid := range deleteServiceUuids {
		buf.WriteString(fmt.Sprintf("\n\t%s", uuid))
	}

	buf.WriteString(fmt.Sprintf("\nStartingServices: %d", len(startingServiceUuids)))
	for _, uuid := range startingServiceUuids {
		buf.WriteString(fmt.Sprintf("\n\t%s", uuid))
	}
	log.Debugf(buf.String() + "\n")

	s.lock.Unlock()
}

func (s *Scheduler) ExecuteService(serv service.ServiceInterface) error {
	// Pass channel because scheduler need to update service's status.
	switch serv.Version() {
	case service.SERVICE_VERSION_V1:
		se := executor.NewServiceExecutor(*serv.(*service.ServiceV1), s.updateChan)
		return se.Execute()
	case service.SERVICE_VERSION_V2:
		se := executor.NewServiceExecutorV2(*serv.(*service.ServiceV2), s.updateChan)
		return se.Execute()
	}

	return fmt.Errorf("not supported service version(%s)", serv.Version())
}

func (s *Scheduler) RecvNotifyServiceStatus() {
	// If you want to stop. close(s.ch).
	for update := range s.updateChan {
		s.notifyUpdateChan <- update
	}
}

func (s *Scheduler) UpdateServiceStatus(update service.ServiceUpdateInterface) {
	serviceStatus := service.ServiceStatusProcessing
	if update.GetStatus() == service.StepStatusFail {
		serviceStatus = service.ServiceStatusFailed
	} else {
		if update.GetStepCount() == update.GetSequence()+1 {
			if update.GetStatus() == service.StepStatusSuccess {
				serviceStatus = service.ServiceStatusSuccess
			}
		}
	}
	s.lock.Lock()
	prevStatus, ok := s.servicesStatusMap[update.GetId()]
	if ok {
		if prevStatus < serviceStatus {
			s.servicesStatusMap[update.GetId()] = service.ServiceStatus(serviceStatus)
		}
	}
	s.lock.Unlock()
}

func (s *Scheduler) NotifyServiceUpdate() <-chan service.ServiceUpdateInterface {
	return s.notifyUpdateChan
}

func (s *Scheduler) CleanupRemainingServices() map[string]service.ServiceStatus {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.servicesStatusMap) == 0 {
		return nil
	}

	services := make(map[string]service.ServiceStatus)

	for uuid, status := range s.servicesStatusMap {
		if status == service.ServiceStatusSuccess || status == service.ServiceStatusFailed {
			delete(s.servicesStatusMap, uuid)
			continue
		}
		services[uuid] = status
	}

	return services
}
