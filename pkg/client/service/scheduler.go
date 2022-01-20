package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

const defaultPeriodicServiceInterval = 5 // * time.Second

type ServiceCheckedFlag int32

const (
	ServiceCheckedFlagCreated = iota
	ServiceCheckedFlagUpdated
	ServiceCheckedFlagSent
	ServiceCheckedFlagDone
)

type ServiceChecked struct {
	flag    ServiceCheckedFlag // 0: created, 1: updated, 2: sent, 3: done
	service Service
}

type ServiceScheduler struct {
	scheduler  *gocron.Scheduler // scheduling services
	services   map[string]*ServiceChecked
	updateChan chan Service // this channel receives service's status
	lock       sync.RWMutex
}

func NewScheduler() *ServiceScheduler {
	return &ServiceScheduler{scheduler: gocron.NewScheduler(time.UTC), services: make(map[string]*ServiceChecked), updateChan: make(chan Service)}
}

func (s *ServiceScheduler) Start() error {
	if s.updateChan == nil || s.services == nil {
		return fmt.Errorf("scheduler don't have channel")
	}

	go s.RecvNotifyServiceStatus()

	s.scheduler.StartAsync()

	return nil
}

func (s *ServiceScheduler) RegisterServices(services map[string]*Service) {
	for _, serv := range services {
		s.lock.Lock()
		s.services[serv.Id] = &ServiceChecked{service: *serv}
		s.lock.Unlock()

		// Create and Execute(goroutine) Service.
		// Pass channel because scheduler need to update service's status.
		go s.ExecuteService(serv)
	}
}

func (s *ServiceScheduler) GetDeleteServicesUpdated() map[string]ServiceChecked {
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
			delete(s.services, id)
		}
	}

	return res
}

func (s *ServiceScheduler) RepairUpdateFailedServices(services map[string]ServiceChecked) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) == 0 {
		return
	}

	for id, v := range services {
		serv, ok := s.services[id]
		if ok {
			if serv.flag == ServiceCheckedFlagSent {
				serv.flag = v.flag
			}
		} else {
			s.services[id] = &v
		}
	}
}

func (s *ServiceScheduler) DeleteDuplicatedServices(services map[string]*Service) map[string]*Service {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) == 0 {
		return services
	}

	for id, _ := range services {
		if _, ok := s.services[id]; ok {
			delete(services, id)
		}
	}

	return services
}

func (s *ServiceScheduler) ExecuteService(serv *Service) error {
	se := NewServiceExecutor(*serv, s.updateChan)

	if serv.ExecType == ServiceExecTypeImmediate {
		return se.Execute()
	} else {
		s.scheduler.Every(defaultPeriodicServiceInterval).Second().Do(se.Execute())
	}
	return nil
}

func (s *ServiceScheduler) RecvNotifyServiceStatus() {
	// If you want to stop. close(s.ch).
	for service := range s.updateChan {
		s.lock.Lock()
		s.services[service.Id].service = service
		if service.Status == ServiceStatusSuccess || service.Status == ServiceStatusFailed {
			s.services[service.Id].flag = ServiceCheckedFlagDone
		} else {
			s.services[service.Id].flag = ServiceCheckedFlagUpdated
		}
		s.lock.Unlock()
	}
}
