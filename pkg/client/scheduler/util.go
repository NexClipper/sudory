package scheduler

import (
	"github.com/NexClipper/sudory/pkg/client/service"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

func ServiceListServerToClient2(server []servicev1.HttpRspClientSideService) map[string]*service.Service {
	client := make(map[string]*service.Service)
	for _, v := range server {
		serv := &service.Service{
			Id:         v.Uuid,
			Name:       *v.Name,
			ClusterId:  *v.ClusterUuid,
			ServerData: v,
		}
		for i, s := range v.Steps {
			serv.Steps = append(serv.Steps, &service.Step{
				Id:       i,
				ParentId: serv.Id,
				Command:  &service.StepCommand{Method: *s.Method, Args: s.Args}})
		}
		client[v.Uuid] = serv
	}

	return client
}

func ServiceListClientToServer2(client map[string]ServiceChecked) []servicev1.HttpReqClientSideService {
	var server []servicev1.HttpReqClientSideService

	if client == nil {
		return server
	}

	for _, v := range client {
		serv := servicev1.HttpReqClientSideService{ServiceAndSteps: servicev1.ServiceAndSteps{Service: v.service.ServerData.Service, Steps: v.service.ServerData.Steps}}
		switch v.service.Status {
		case service.ServiceStatusPreparing, service.ServiceStatusStart, service.ServiceStatusProcessing:
			serv.Service.Status = newist.Int32(int32(servicev1.StatusProcessing))
		case service.ServiceStatusSuccess:
			serv.Service.Status = newist.Int32(int32(servicev1.StatusSuccess))
		case service.ServiceStatusFailed:
			serv.Service.Status = newist.Int32(int32(servicev1.StatusFail))
		}

		if v.service.Result.Body != "" {
			serv.Service.Result = newist.String(v.service.Result.Body)
		}
		if v.service.Result.Err != nil {
			serv.Service.Result = newist.String(v.service.Result.Err.Error())
		}

		for i, s := range v.service.Steps {
			serv.Steps[i].Status = newist.Int32(int32(s.Status))
			switch s.Status {
			case service.StepStatusPreparing, service.StepStatusProcessing:
				serv.Steps[i].Status = newist.Int32(int32(servicev1.StatusProcessing))
			case service.StepStatusSuccess:
				serv.Steps[i].Status = newist.Int32(int32(servicev1.StatusSuccess))
			case service.StepStatusFail:
				serv.Steps[i].Status = newist.Int32(int32(servicev1.StatusFail))
			}
		}

		server = append(server, serv)
	}

	return server
}
