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
	updateChan        chan service.UpdateServiceStep // this channel receives service's status
	notifyUpdateChan  chan service.UpdateServiceStep
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		servicesStatusMap: make(map[string]service.ServiceStatus),
		maxProcessLimit:   defaultMaxProcessLimit,
		updateChan:        make(chan service.UpdateServiceStep),
		notifyUpdateChan:  make(chan service.UpdateServiceStep)}
}

func (s *Scheduler) Start() error {
	if s.updateChan == nil || s.servicesStatusMap == nil {
		return fmt.Errorf("scheduler don't have channel")
	}

	go s.RecvNotifyServiceStatus()

	return nil
}

func (s *Scheduler) RegisterServices(services map[string]*service.Service) {
	// 1. already existing services drop
	var startingList []*service.Service
	s.lock.Lock()
	for _, service := range services {
		_, ok := s.servicesStatusMap[service.Id]
		if !ok {
			startingList = append(startingList, service)
		}
	}
	sort.Slice(startingList, func(i, j int) bool {
		if startingList[i].Priority > startingList[j].Priority {
			return true
		} else if startingList[i].Priority < startingList[j].Priority {
			return false
		} else {
			return startingList[i].CreatedTime.Before(startingList[j].CreatedTime)
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
			startingServiceUuids = append(startingServiceUuids, serv.Id)
			s.servicesStatusMap[serv.Id] = service.ServiceStatusPreparing
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

func (s *Scheduler) ExecuteService(serv *service.Service) error {
	// Pass channel because scheduler need to update service's status.
	se := executor.NewServiceExecutor(*serv, s.updateChan)

	return se.Execute()
}

func (s *Scheduler) RecvNotifyServiceStatus() {
	// If you want to stop. close(s.ch).
	for update := range s.updateChan {
		s.notifyUpdateChan <- update
	}
}

func (s *Scheduler) UpdateServiceStatus(update service.UpdateServiceStep) {
	serviceStatus := service.ServiceStatusProcessing
	if update.StepCount == update.Sequence+1 {
		if update.Status == service.StepStatusSuccess {
			serviceStatus = service.ServiceStatusSuccess
		} else if update.Status == service.StepStatusFail {
			serviceStatus = service.ServiceStatusFailed
		}
	}
	s.lock.Lock()
	prevStatus, ok := s.servicesStatusMap[update.Uuid]
	if ok {
		if prevStatus < serviceStatus {
			s.servicesStatusMap[update.Uuid] = service.ServiceStatus(serviceStatus)
		}
	}
	s.lock.Unlock()
}

func (s *Scheduler) NotifyServiceUpdate() <-chan service.UpdateServiceStep {
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
