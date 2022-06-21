package service

import (
	"time"

	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
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
	StartTime  time.Time
	UpdateTime time.Time
	EndTime    time.Time
	Status     ServiceStatus
	Steps      []Step
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

type UpdateServiceStep struct {
	Uuid      string
	StepCount int
	Sequence  int
	Status    StepStatus
	Result    string
	Started   time.Time
	Ended     time.Time
}

func ConvertServiceListServerToClient(server []servicev1.HttpRspService_ClientSide) map[string]*Service {
	client := make(map[string]*Service)
	for _, v := range server {
		serv := &Service{
			Id:         v.Uuid,
			Name:       v.Name,
			ClusterId:  v.ClusterUuid,
			ServerData: v,
		}
		for _, s := range v.Steps {
			rf := ""
			if s.ResultFilter != nil {
				rf = *s.ResultFilter
			}
			serv.Steps = append(serv.Steps, Step{
				Id:           s.Uuid,
				ParentId:     serv.Id,
				Command:      &StepCommand{Method: s.Method, Args: s.Args},
				ResultFilter: rf,
			})
		}
		client[v.Uuid] = serv
	}

	return client
}

func ConvertServiceStepUpdateClientToServer(client UpdateServiceStep) *servicev2.HttpReq_ClientServiceUpdate {
	server := &servicev2.HttpReq_ClientServiceUpdate{
		Uuid:     client.Uuid,
		Sequence: client.Sequence,
		// Status:client.Status,
		Result:  client.Result,
		Started: client.Started,
		Ended:   client.Ended,
	}

	switch client.Status {
	case StepStatusPreparing, StepStatusProcessing:
		server.Status = servicev2.StepStatusProcessing
	case StepStatusSuccess:
		server.Status = servicev2.StepStatusSuccess
	case StepStatusFail:
		server.Status = servicev2.StepStatusFail
	}

	return server
}
