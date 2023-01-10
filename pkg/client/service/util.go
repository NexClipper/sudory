package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	servicev4 "github.com/NexClipper/sudory/pkg/server/model/service/v4"
)

type FailedConvertService struct {
	Data servicev4.HttpRsp_ClientServicePolling
	Err  error
}

func ConvertServiceListServerToClient(server []servicev4.HttpRsp_ClientServicePolling) (map[string]ServiceInterface, []FailedConvertService) {
	client := make(map[string]ServiceInterface)
	var failed []FailedConvertService
	for _, v := range server {
		switch v.Version {
		case "v3":
			serv := &ServiceV1{
				Id:          v.V3.Uuid,
				Name:        v.V3.Name,
				ClusterId:   v.V3.ClusterUuid,
				Priority:    int(v.V3.Priority),
				CreatedTime: v.V3.Created,
			}

			if len(v.V3.Steps) <= 0 {
				log.Warnf("service steps is empty: service_uuid: %s\n", v.V3.Uuid)
				continue
			}

			for _, s := range v.V3.Steps {
				serv.Steps = append(serv.Steps, Step{
					Id:           s.Uuid,
					ParentId:     serv.Id,
					Command:      &StepCommand{Method: s.Method, Args: s.Args},
					ResultFilter: s.ResultFilter.String,
				})
			}
			client[v.V3.Uuid] = serv
		case "v4":
			serv := &ServiceV2{
				Id:          v.V4.Uuid,
				Name:        v.V4.Name,
				ClusterId:   v.V4.ClusterUuid,
				Priority:    int(v.V4.Priority),
				CreatedTime: v.V4.Created,
			}

			var flow Flow
			if err := json.Unmarshal([]byte(v.V4.Flow), &flow); err != nil {
				log.Warnf("failed to convert service(%s)\n", v.V4.Uuid)
				failed = append(failed, FailedConvertService{Data: v, Err: err})
				continue
			}

			if len(flow) <= 0 {
				log.Warnf("service steps is empty: service_uuid: %s\n", v.V4.Uuid)
				continue
			}

			for _, flowstep := range flow {
				if err := flowstep.Inputs.FindReplaceDeferredInputsFrom(v.V4.Inputs); err != nil {
					log.Warnf("failed to convert service(%s)\n", v.V4.Uuid)
					failed = append(failed, FailedConvertService{Data: v, Err: err})
					continue
				}
			}

			serv.Flow = flow

			client[v.V4.Uuid] = serv
		default:
			err := fmt.Errorf("unknown service version(%s)", v.Version)
			log.Warnf("failed to convert service(%s): %s\n", v.V4.Uuid, err.Error())
			failed = append(failed, FailedConvertService{Data: v, Err: err})
			continue
		}
	}

	return client, failed
}

func ConvertServiceStepUpdateClientToServer(client ServiceUpdateInterface) *servicev4.HttpReq_ClientServiceUpdate {
	var server *servicev4.HttpReq_ClientServiceUpdate
	switch client.Version().GetMatchServerServiceVersion() {
	case "v3":
		update := client.(*UpdateServiceV1)
		server = &servicev4.HttpReq_ClientServiceUpdate{
			Version: client.Version().GetMatchServerServiceVersion(),
			V3: servicev3.HttpReq_ClientServiceUpdate{
				Uuid:     update.Uuid,
				Sequence: update.Sequence,
				Result:   update.Result,
				Started:  update.Started,
				Ended:    update.Ended,
			},
		}

		switch update.Status {
		case StepStatusPreparing, StepStatusProcessing:
			server.V3.Status = servicev3.StepStatusProcessing
		case StepStatusSuccess:
			server.V3.Status = servicev3.StepStatusSuccess
		case StepStatusFail:
			server.V3.Status = servicev3.StepStatusFail
		}
	case "v4":
		update := client.(*UpdateServiceV2)
		server = &servicev4.HttpReq_ClientServiceUpdate{
			Version: client.Version().GetMatchServerServiceVersion(),
			V4: servicev4.HttpReq_ClientServiceUpdate_multistep{
				Uuid:     update.Id,
				Sequence: update.Sequence,
				Result:   update.Result,
				Started:  update.Started,
				Ended:    update.Ended,
			},
		}

		switch update.Status {
		case StepStatusPreparing, StepStatusProcessing:
			server.V4.Status = servicev4.StepStatusProcessing
		case StepStatusSuccess:
			server.V4.Status = servicev4.StepStatusSucceeded
		case StepStatusFail:
			server.V4.Status = servicev4.StepStatusFailed
		}
	}

	return server
}

func CreateUpdateService(version Version, id string, seqCount, seq int, status StepStatus, result string, start, end time.Time) ServiceUpdateInterface {
	switch version {
	case SERVICE_VERSION_V1:
		return &UpdateServiceV1{
			Uuid:      id,
			StepCount: seqCount,
			Sequence:  seq,
			Status:    status,
			Result:    result,
			Started:   start,
			Ended:     end,
		}
	case SERVICE_VERSION_V2:
		return &UpdateServiceV2{
			Id:        id,
			StepCount: seqCount,
			Sequence:  seq,
			Status:    status,
			Result:    result,
			Started:   start,
			Ended:     end,
		}
	}
	return nil
}
