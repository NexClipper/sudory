package service

import (
	"time"
)

type ServiceStatus int32

const (
	ServiceStatusPreparing ServiceStatus = iota + 1
	ServiceStatusStart
	ServiceStatusProcessing
	ServiceStatusSuccess
	ServiceStatusFailed
)

func (s ServiceStatus) String() string {
	switch s {
	case ServiceStatusPreparing:
		return "ServiceStatusPreparing"
	case ServiceStatusStart:
		return "ServiceStatusStart"
	case ServiceStatusProcessing:
		return "ServiceStatusProcessing"
	case ServiceStatusSuccess:
		return "ServiceStatusSuccess"
	case ServiceStatusFailed:
		return "ServiceStatusFailed"
	default:
		return "ServiceStatusUnknown"
	}
}

type ServiceV1 struct {
	Id          string
	Name        string
	ClusterId   string
	Priority    int
	CreatedTime time.Time
	StartTime   time.Time
	UpdateTime  time.Time
	EndTime     time.Time
	Status      ServiceStatus
	Steps       []Step
	Result      Result
}

func (s *ServiceV1) Version() Version {
	return SERVICE_VERSION_V1
}

func (s *ServiceV1) GetId() string {
	return s.Id
}

func (s *ServiceV1) GetPriority() int {
	return s.Priority
}

func (s *ServiceV1) GetCreatedTime() time.Time {
	return s.CreatedTime
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

type StepCommand struct {
	Method string
	Args   map[string]interface{}
}

type UpdateServiceV1 struct {
	Uuid      string
	StepCount int
	Sequence  int
	Status    StepStatus
	Result    string
	Started   time.Time
	Ended     time.Time
}

func (s *UpdateServiceV1) Version() Version {
	return SERVICE_VERSION_V1
}

func (s *UpdateServiceV1) GetId() string {
	return s.Uuid
}

func (s *UpdateServiceV1) GetStepCount() int {
	return s.StepCount
}

func (s *UpdateServiceV1) GetSequence() int {
	return s.Sequence
}

func (s *UpdateServiceV1) GetStatus() StepStatus {
	return s.Status
}

// func ConvertServiceListServerToClient(server []servicev3.HttpRsp_ClientServicePolling) map[string]*ServiceV1 {
// 	client := make(map[string]*ServiceV1)
// 	for _, v := range server {
// 		serv := &ServiceV1{
// 			Id:          v.Uuid,
// 			Name:        v.Name,
// 			ClusterId:   v.ClusterUuid,
// 			Priority:    int(v.Priority),
// 			CreatedTime: v.Created,
// 		}

// 		if len(v.Steps) <= 0 {
// 			log.Warnf("service steps is empty: service_uuid: %s\n", v.Uuid)
// 			continue
// 		}

// 		for _, s := range v.Steps {
// 			serv.Steps = append(serv.Steps, Step{
// 				Id:           s.Uuid,
// 				ParentId:     serv.Id,
// 				Command:      &StepCommand{Method: s.Method, Args: s.Args},
// 				ResultFilter: s.ResultFilter.String,
// 			})
// 		}
// 		client[v.Uuid] = serv
// 	}

// 	return client
// }

// func ConvertServiceStepUpdateClientToServer(client UpdateServiceV1) *servicev3.HttpReq_ClientServiceUpdate {
// 	server := &servicev3.HttpReq_ClientServiceUpdate{
// 		Uuid:     client.Uuid,
// 		Sequence: client.Sequence,
// 		// Status:client.Status,
// 		Result:  client.Result,
// 		Started: client.Started,
// 		Ended:   client.Ended,
// 	}

// 	switch client.Status {
// 	case StepStatusPreparing, StepStatusProcessing:
// 		server.Status = servicev3.StepStatusProcessing
// 	case StepStatusSuccess:
// 		server.Status = servicev3.StepStatusSuccess
// 	case StepStatusFail:
// 		server.Status = servicev3.StepStatusFail
// 	}

// 	return server
// }
