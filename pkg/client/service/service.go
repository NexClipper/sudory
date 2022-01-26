package service

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/newist"
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
	serverData servicev1.HttpRspClientSideService
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
	Args   map[string]string
}

type Result struct {
	body string
	err  error
}

type Step struct {
	Id        int
	ParentId  string
	Command   *StepCommand
	StartTime time.Time
	EndTime   time.Time
	Status    StepStatus
	Result    Result
}

func ServiceListServerToClient(server []servicev1.HttpRspClientSideService) map[string]*Service {
	client := make(map[string]*Service)
	for _, v := range server {
		serv := &Service{
			Id:         v.Uuid,
			Name:       v.Name,
			ClusterId:  v.ClusterUuid,
			serverData: v,
		}
		for i, s := range v.Steps {
			serv.Steps = append(serv.Steps, &Step{
				Id:       i,
				ParentId: serv.Id,
				Command:  &StepCommand{Method: *s.Method, Args: s.Args}})
		}
		client[v.Uuid] = serv
	}

	return client
}

func ServiceListClientToServer(client map[string]ServiceChecked) []servicev1.HttpReqClientSideService {
	var server []servicev1.HttpReqClientSideService

	if client == nil {
		return server
	}

	for _, v := range client {
		serv := servicev1.HttpReqClientSideService{Service: v.service.serverData.Service, Steps: v.service.serverData.Steps}
		switch v.service.Status {
		case ServiceStatusPreparing, ServiceStatusStart, ServiceStatusProcessing:
			serv.Service.Status = newist.Int32(int32(servicev1.StatusProcessing))
		case ServiceStatusSuccess:
			serv.Service.Status = newist.Int32(int32(servicev1.StatusSuccess))
		case ServiceStatusFailed:
			serv.Service.Status = newist.Int32(int32(servicev1.StatusFail))
		}

		if v.service.Result.body != "" {
			serv.Service.Result = newist.String(v.service.Result.body)
		}
		if v.service.Result.err != nil {
			serv.Service.Result = newist.String(v.service.Result.err.Error())
		}

		for i, s := range v.service.Steps {
			serv.Steps[i].Status = newist.Int32(int32(s.Status))
		}

		server = append(server, serv)
	}

	return server
}
