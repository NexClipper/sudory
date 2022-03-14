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
	ServiceCheckedFlagSent
	ServiceCheckedFlagDone
)

type ServiceChecked struct {
	flag    ServiceCheckedFlag // 0: created, 1: updated, 2: sent, 3: done
	service service.Service
}

type Scheduler struct {
	services   map[string]*ServiceChecked
	updateChan chan service.Service // this channel receives service's status
	lock       sync.RWMutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{services: make(map[string]*ServiceChecked), updateChan: make(chan service.Service)}
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
		s.services[serv.Id].service = serv
		if serv.Status == service.ServiceStatusSuccess || serv.Status == service.ServiceStatusFailed {
			s.services[serv.Id].flag = ServiceCheckedFlagDone
		} else {
			s.services[serv.Id].flag = ServiceCheckedFlagUpdated
		}
		s.lock.Unlock()
	}
}

// get services with updated status
func (s *Scheduler) GetServicesWithUpdatedDoneFlag() map[string]ServiceChecked {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) <= 0 {
		return nil
	}

	res := make(map[string]ServiceChecked)
	for id, serv := range s.services {
		if serv.flag == ServiceCheckedFlagUpdated {
			res[id] = ServiceChecked{flag: serv.flag, service: serv.service}
			serv.flag = ServiceCheckedFlagSent
		} else if serv.flag == ServiceCheckedFlagDone {
			res[id] = ServiceChecked{flag: serv.flag, service: serv.service}
		}
	}

	return res
}

func (s *Scheduler) DeleteServicesWithDoneFlag(services map[string]ServiceChecked) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) <= 0 {
		return
	}

	for id, v := range services {
		if v.flag == ServiceCheckedFlagDone {
			delete(s.services, id)
		}
	}
}

func (s *Scheduler) RollbackServicesWithDoneUpdatedFlag(services map[string]ServiceChecked) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) <= 0 {
		return
	}

	for id, v := range services {
		if v.flag == ServiceCheckedFlagDone {
			s.services[id].flag = v.flag
		} else if v.flag == ServiceCheckedFlagUpdated {
			if s.services[id].flag == ServiceCheckedFlagSent {
				s.services[id].flag = v.flag
			}
		}
	}
}
