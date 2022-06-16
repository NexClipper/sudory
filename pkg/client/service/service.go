package service

import (
	"time"

	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

type ServiceExecType int32

const (
	ServiceExecTypeImmediate = iota
	// ServiceExecTypePeriodic
)

type ServiceStatus int32

const (
	ServiceStatusPreparing = iota + 1
	ServiceStatusStart
	ServiceStatusProcessing
	ServiceStatusSuccess
	ServiceStatusFailed
)

type Service struct {
	Id         string
	Name       string
	ClusterId  string
	ExecType   ServiceExecType
	StartTime  time.Time
	UpdateTime time.Time
	EndTime    time.Time
	Status     ServiceStatus
	Steps      []*Step
	Result     Result
	ServerData servicev1.HttpRspService_ClientSide
}

type StepStatus int32

const (
	StepStatusPreparing = iota + 1
	StepStatusProcessing
	StepStatusSuccess
	StepStatusFail
)

type StepCommand struct {
	Method string
	Args   map[string]interface{}
}

type Result struct {
	Body string
	Err  error
}

type Step struct {
	Id           string
	ParentId     string
	Command      *StepCommand
	StartTime    time.Time
	EndTime      time.Time
	Status       StepStatus
	ResultFilter string
	Result       Result
}
