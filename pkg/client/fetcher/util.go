package fetcher

import (
	"time"

	"github.com/NexClipper/sudory/pkg/client/service"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

func ConvertServiceListServerToClient(server []servicev1.HttpRspService_ClientSide) map[string]*service.Service {
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

func ConvertServiceClientToServer(client *service.Service) *ReqUpdateService {
	if client == nil {
		return nil
	}

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
		case service.StepStatusPreparing, service.StepStatusProcessing:
			st.Status = int32(servicev1.StatusProcessing)
		case service.StepStatusSuccess:
			st.Status = int32(servicev1.StatusSuccess)
		case service.StepStatusFail:
			st.Status = int32(servicev1.StatusFail)
		}

		server.Steps = append(server.Steps, st)
	}

	return server
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
