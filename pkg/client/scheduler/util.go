package scheduler

import (
	"github.com/NexClipper/sudory/pkg/client/service"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

func ServiceListServerToClient(server []servicev1.HttpRspService_ClientSide) map[string]*service.Service {
	client := make(map[string]*service.Service)
	for _, v := range server {
		serv := &service.Service{
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
			serv.Steps = append(serv.Steps, service.Step{
				Id:           s.Uuid,
				ParentId:     serv.Id,
				Command:      &service.StepCommand{Method: s.Method, Args: s.Args},
				ResultFilter: rf,
			})
		}
		client[v.Uuid] = serv
	}

	return client
}

func ServiceClientToServer(client *ServiceChecked) *servicev1.HttpReq_ServiceUpdate_ClientSide {
	if client == nil {
		return nil
	}

	server := &servicev1.HttpReq_ServiceUpdate_ClientSide{
		UuidMeta: metav1.UuidMeta{Uuid: client.service.Id},
	}

	if client.service.Result.Body != "" {
		server.Result = &client.service.Result.Body
	}
	if client.service.Result.Err != nil {
		err := client.service.Result.Err.Error()
		server.Result = &err
	}

	for _, s := range client.service.Steps {
		st := servicev1.HttpReq_ServiceUpdate_Step_ClientSide{
			UuidMeta: metav1.UuidMeta{Uuid: s.Id},
			Started:  newist.Time(s.StartTime),
			Ended:    newist.Time(s.EndTime),
		}

		switch s.Status {
		case service.StepStatusPreparing, service.StepStatusProcessing:
			st.Status = newist.Int32(int32(servicev1.StatusProcessing))
		case service.StepStatusSuccess:
			st.Status = newist.Int32(int32(servicev1.StatusSuccess))
		case service.StepStatusFail:
			st.Status = newist.Int32(int32(servicev1.StatusFail))
		}

		server.Steps = append(server.Steps, st)
	}

	return server
}

func ServiceListClientToServer(client map[string]*ServiceChecked) []*servicev1.HttpReq_ServiceUpdate_ClientSide {
	if client == nil {
		return nil
	}

	server := make([]*servicev1.HttpReq_ServiceUpdate_ClientSide, 0, len(client))
	for _, v := range client {
		serv := &servicev1.HttpReq_ServiceUpdate_ClientSide{
			UuidMeta: metav1.UuidMeta{Uuid: v.service.Id},
		}

		if v.service.Result.Body != "" {
			serv.Result = &v.service.Result.Body
		}
		if v.service.Result.Err != nil {
			err := v.service.Result.Err.Error()
			serv.Result = &err
		}

		for _, s := range v.service.Steps {
			st := servicev1.HttpReq_ServiceUpdate_Step_ClientSide{UuidMeta: metav1.UuidMeta{Uuid: s.Id}}

			switch s.Status {
			case service.StepStatusPreparing, service.StepStatusProcessing:
				st.Status = newist.Int32(int32(servicev1.StatusProcessing))
			case service.StepStatusSuccess:
				st.Status = newist.Int32(int32(servicev1.StatusSuccess))
			case service.StepStatusFail:
				st.Status = newist.Int32(int32(servicev1.StatusFail))
			}
			serv.Steps = append(serv.Steps, st)
		}

		server = append(server, serv)
	}

	return server
}
