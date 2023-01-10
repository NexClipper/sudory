package service

import (
	"time"
)

type ServiceV2 struct {
	Id          string
	Name        string
	ClusterId   string
	Priority    int
	CreatedTime time.Time
	StartTime   time.Time
	UpdateTime  time.Time
	EndTime     time.Time
	Status      ServiceStatus
	Flow        Flow
	// Inputs      map[string]interface{}
	Result      Result
}

func (s *ServiceV2) Version() Version {
	return SERVICE_VERSION_V2
}

func (s *ServiceV2) GetId() string {
	return s.Id
}

func (s *ServiceV2) GetPriority() int {
	return s.Priority
}

func (s *ServiceV2) GetCreatedTime() time.Time {
	return s.CreatedTime
}

type UpdateServiceV2 struct {
	Id        string
	StepCount int
	Sequence  int
	Status    StepStatus
	Result    string
	Started   time.Time
	Ended     time.Time
}

func (s *UpdateServiceV2) Version() Version {
	return SERVICE_VERSION_V2
}

func (s *UpdateServiceV2) GetId() string {
	return s.Id
}

func (s *UpdateServiceV2) GetStepCount() int {
	return s.StepCount
}

func (s *UpdateServiceV2) GetSequence() int {
	return s.Sequence
}

func (s *UpdateServiceV2) GetStatus() StepStatus {
	return s.Status
}
