package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

const defaultPeriodicServiceInterval = 5 // * time.Second

type ServiceScheduler struct {
	scheduler  *gocron.Scheduler // scheduling services
	services   map[string]*Service
	updateChan chan *Service // this channel receives service's status
	lock       sync.RWMutex
}

func NewScheduler() *ServiceScheduler {
	return &ServiceScheduler{scheduler: gocron.NewScheduler(time.UTC), services: make(map[string]*Service), updateChan: make(chan *Service)}
}

func (s *ServiceScheduler) Start() error {
	if s.updateChan == nil || s.services == nil {
		return fmt.Errorf("scheduler don't have channel.")
	}

	go s.RecvNotifyServiceStatus()

	s.scheduler.StartAsync()

	return nil
}

func (s *ServiceScheduler) RegisterServices(services map[string]*Service) {
	for _, serv := range services {
		s.lock.Lock()
		s.services[serv.Id] = serv
		s.lock.Unlock()

		// Create and Execute(goroutine) Service.
		// Pass channel because scheduler need to update service's status.
		go s.ExecuteService(serv)
	}
}

func (s *ServiceScheduler) GetServices() map[string]*Service {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if len(s.services) <= 0 {
		return nil
	}

	res := make(map[string]*Service)
	for id, serv := range s.services {
		res[id] = serv
	}

	return res
}

func (s *ServiceScheduler) RemoveServices(services map[string]*Service) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.services) <= 0 {
		return
	}

	for id, _ := range services {
		delete(s.services, id)
	}
}

func (s *ServiceScheduler) ExecuteService(serv *Service) error {
	if serv.Type == ServiceTypeAtOnce {
		return serv.Execute(s.updateChan)
	} else {
		s.scheduler.Every(defaultPeriodicServiceInterval).Second().Do(serv.Execute, s.updateChan)
	}
	return nil
}

func (s *ServiceScheduler) RecvNotifyServiceStatus() {
	// If you want to stop. close(s.ch).
	for serv := range s.updateChan {
		s.lock.Lock()
		s.services[serv.Id] = serv
		s.lock.Unlock()
	}
}
