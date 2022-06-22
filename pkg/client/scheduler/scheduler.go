package scheduler

import (
	"fmt"
	"sync"

	"github.com/NexClipper/sudory/pkg/client/executor"
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

	// 2. if existing service's status is ServiceStatusSuccess or ServiceStatusFailed, delete in statusMap
	for uuid, status := range s.servicesStatusMap {
		if status == service.ServiceStatusSuccess || status == service.ServiceStatusFailed {
			delete(s.servicesStatusMap, uuid)
		}
	}

	// 3. maxProcessLimit - len(statusMap) is number starting now
	remain := s.maxProcessLimit - len(s.servicesStatusMap)

	for _, serv := range startingList {
		if remain > 0 {
			s.servicesStatusMap[serv.Id] = service.ServiceStatusPreparing
			// create and execute(goroutine) service.
			go s.ExecuteService(serv)
			remain--
		} else {
			break
		}
	}
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
		serviceStatus := service.ServiceStatusProcessing
		if update.StepCount == update.Sequence+1 {
			if update.Status == service.StepStatusSuccess {
				serviceStatus = service.ServiceStatusSuccess
			} else if update.Status == service.StepStatusFail {
				serviceStatus = service.ServiceStatusFailed
			}
		}
		s.lock.Lock()
		s.servicesStatusMap[update.Uuid] = service.ServiceStatus(serviceStatus)
		s.lock.Unlock()
		s.notifyUpdateChan <- update
	}
}

func (s *Scheduler) NotifyServiceUpdate() <-chan service.UpdateServiceStep {
	return s.notifyUpdateChan
}
