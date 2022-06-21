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

type ReqUpdateService struct {
	Uuid   string                   `json:"uuid"`
	Result string                   `json:"result,omitempty"`
	Steps  []*ReqUpdateService_Step `json:"steps,omitempty"`
}

type ReqUpdateService_Step struct {
	Uuid    string    `json:"uuid"`
	Status  int32     `json:"status,omitempty"`
	Started time.Time `json:"started,omitempty"`
	Ended   time.Time `json:"ended,omitempty"`
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

func ConvertServiceClientToServer(client Service) *ReqUpdateService {
	server := &ReqUpdateService{
		Uuid: client.Id,
	}

	if client.Result.Body != "" {
		server.Result = client.Result.Body
	}
	if client.Result.Err != nil {
		err := client.Result.Err.Error()
		server.Result = err
	}

	for _, s := range client.Steps {
		st := &ReqUpdateService_Step{
			Uuid:    s.Id,
			Started: s.StartTime,
			Ended:   s.EndTime,
		}

		switch s.Status {
		case StepStatusPreparing, StepStatusProcessing:
			st.Status = int32(servicev1.StatusProcessing)
		case StepStatusSuccess:
			st.Status = int32(servicev1.StatusSuccess)
		case StepStatusFail:
			st.Status = int32(servicev1.StatusFail)
		}

		server.Steps = append(server.Steps, st)
	}

	return server
}
