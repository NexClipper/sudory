package scheduler

import (
	"fmt"
	"sync"

	"github.com/NexClipper/sudory/pkg/client/executor"
	"github.com/NexClipper/sudory/pkg/client/service"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
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
	notifyUpdateChan chan *servicev1.HttpReq_ServiceUpdate_ClientSide
	lock             sync.RWMutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		services:         make(map[string]*ServiceChecked),
		updateChan:       make(chan service.Service),
		notifyUpdateChan: make(chan *servicev1.HttpReq_ServiceUpdate_ClientSide)}
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
		var send *servicev1.HttpReq_ServiceUpdate_ClientSide
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

		send = ServiceClientToServer(sc)
		s.lock.Unlock()
		s.notifyUpdateChan <- send
	}
}

func (s *Scheduler) NotifyServiceUpdate() <-chan *servicev1.HttpReq_ServiceUpdate_ClientSide {
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

func (s *Scheduler) ChangeServiceFlagFailedToSend(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	sc, ok := s.services[id]
	if ok {
		sc.flag = ServiceCheckedFlagFailedToSend
	}
}
