package scheduler

import (
	"fmt"
	"sync"

	"github.com/NexClipper/sudory/pkg/client/executor"
	"github.com/NexClipper/sudory/pkg/client/service"
)

type ServiceCheckedFlag int32

const (
	ServiceCheckedFlagCreated ServiceCheckedFlag = iota
	ServiceCheckedFlagUpdated
	ServiceCheckedFlagDone
	ServiceCheckedFlagFailedToSend
)

type ServiceChecked struct {
	flag    ServiceCheckedFlag // 0: created, 1: updated, 2: done, 3: failed_to_send
	service service.Service
}

func (sc *ServiceChecked) GetServiceId() string {
	return sc.service.Id
}

type Scheduler struct {
	services         map[string]*ServiceChecked
	updateChan       chan service.Service // this channel receives service's status
	notifyUpdateChan chan *ServiceChecked
	lock             sync.RWMutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		services:         make(map[string]*ServiceChecked),
		updateChan:       make(chan service.Service),
		notifyUpdateChan: make(chan *ServiceChecked)}
}

func (s *Scheduler) Start() error {
	if s.updateChan == nil || s.services == nil {
		return fmt.Errorf("scheduler don't have channel")
	}

	go s.RecvNotifyServiceStatus()

	return nil
}

func (s *Scheduler) RegisterServices(services map[string]*service.Service) {
	for _, serv := range services {
		s.lock.Lock()

		_, ok := s.services[serv.Id]
		if ok {
			// drop duplicated service
			s.lock.Unlock()
			continue
		}

		s.services[serv.Id] = &ServiceChecked{service: *serv}
		s.lock.Unlock()

		// Create and Execute(goroutine) Service.
		go s.ExecuteService(serv)
	}
}

func (s *Scheduler) ExecuteService(serv *service.Service) error {
	// Pass channel because scheduler need to update service's status.
	se := executor.NewServiceExecutor(*serv, s.updateChan)

	if serv.ExecType == service.ServiceExecTypeImmediate {
		return se.Execute()
	} else {
		// TODO
	}
	return nil
}

func (s *Scheduler) RecvNotifyServiceStatus() {
	// If you want to stop. close(s.ch).
	for serv := range s.updateChan {
		s.lock.Lock()
		sc, ok := s.services[serv.Id]
		if !ok {
			s.lock.Unlock()
			continue
		}

		sc.service = serv
		if serv.Status == service.ServiceStatusSuccess || serv.Status == service.ServiceStatusFailed {
			sc.flag = ServiceCheckedFlagDone
		} else {
			sc.flag = ServiceCheckedFlagUpdated
		}
		s.lock.Unlock()
		s.notifyUpdateChan <- sc
	}
}

func (s *Scheduler) NotifyServiceUpdate() <-chan *ServiceChecked {
	return s.notifyUpdateChan
}

// get services with failed to send
func (s *Scheduler) GetServicesWithFailedToSendFlag() []*ServiceChecked {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) == 0 {
		return nil
	}

	var res []*ServiceChecked
	for _, serv := range s.services {
		if serv.flag == ServiceCheckedFlagFailedToSend {
			res = append(res, serv)
		}
	}

	return res
}

func (s *Scheduler) DeleteServiceWithFlag(id string, flag ServiceCheckedFlag) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if id == "" {
		return
	}

	sc, ok := s.services[id]
	if !ok {
		return
	}

	if sc.flag == flag {
		delete(s.services, sc.service.Id)
	}
}

func (s *Scheduler) ChangeServiceFlagFailedToSend(sc *ServiceChecked) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if sc == nil {
		return
	}

	if sc.flag == ServiceCheckedFlagDone {
		sc.flag = ServiceCheckedFlagFailedToSend
		s.services[sc.service.Id] = sc
	}
}
