package executor

import (
	"encoding/json"
	"time"

	"github.com/NexClipper/sudory/pkg/client/jq"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
)

type ServiceExecutor struct {
	service       service.Service
	updateChannel chan<- service.UpdateServiceStep
}

func NewServiceExecutor(service service.Service, updateChannel chan<- service.UpdateServiceStep) *ServiceExecutor {
	return &ServiceExecutor{service: service, updateChannel: updateChannel}
}

func (se *ServiceExecutor) Execute() (err error) {
	var result service.Result

	defer func() {
		se.service.EndTime = time.Now()
		if err != nil {
			log.Errorf("Failed to execute service: service_uuid: %s, error: %s\n", se.service.Id, err.Error())
			se.service.Status = service.ServiceStatusFailed
			se.service.Result.Err = err
			// se.SendStatusUpdate()
		} else {
			se.service.Status = service.ServiceStatusSuccess
			se.service.Result = result
			// se.SendStatusUpdate()
		}
	}()

	se.service.StartTime = time.Now()
	se.service.Status = service.ServiceStatusProcessing
	//se.SendStatusUpdate()
	for i, step := range se.service.Steps {
		var te *StepExecutor

		// update execute result to service scheduler through returnChannel.
		se.service.Steps[i].Status = service.StepStatusProcessing
		se.service.Steps[i].StartTime = time.Now()
		se.SendServiceStatusUpdate(i, service.StepStatusProcessing, "", se.service.Steps[i].StartTime, se.service.Steps[i].EndTime)

		te, err = NewStepExecutor(step)
		if err != nil {
			se.service.Steps[i].Status = service.StepStatusFail
			se.service.Steps[i].EndTime = time.Now()
			se.SendServiceStatusUpdate(i, se.service.Steps[i].Status, err.Error(), se.service.Steps[i].StartTime, se.service.Steps[i].EndTime)
			return err
		}

		result = te.Execute()
		if err = result.Err; err != nil {
			se.service.Steps[i].Status = service.StepStatusFail
			se.service.Steps[i].EndTime = time.Now()
			se.SendServiceStatusUpdate(i, se.service.Steps[i].Status, err.Error(), se.service.Steps[i].StartTime, se.service.Steps[i].EndTime)
			return err
		}

		// update execute result to service scheduler through returnChannel.
		se.service.Steps[i].Status = service.StepStatusSuccess
		se.service.Steps[i].EndTime = time.Now()
		//se.SendStatusUpdate()
		se.SendServiceStatusUpdate(i, se.service.Steps[i].Status, result.Body, se.service.Steps[i].StartTime, se.service.Steps[i].EndTime)
	}

	return nil
}

// func (se *ServiceExecutor) SendStatusUpdate() {
// 	if se.updateChannel != nil {
// 		se.updateChannel <- se.service
// 	}
// }

func (se *ServiceExecutor) SendServiceStatusUpdate(seq int, status service.StepStatus, result string, st, et time.Time) {
	if se.updateChannel != nil {
		update := service.UpdateServiceStep{
			Uuid:      se.service.Id,
			StepCount: len(se.service.Steps),
			Sequence:  seq,
			Status:    status,
			Result:    result,
			Started:   st,
			Ended:     et,
		}

		se.updateChannel <- update
	}
}

type StepExecutor struct {
	commander Commander
	step      service.Step
}

func NewStepExecutor(step service.Step) (*StepExecutor, error) {
	commander, err := NewCommander(step.Command)
	if err != nil {
		return nil, err
	}

	return &StepExecutor{commander: commander, step: step}, nil
}

func (se *StepExecutor) Execute() service.Result {
	log.Debugf("Prepare to Execute method : %s\n", se.step.Command.Method)
	res, err := se.commander.Run()
	if err != nil {
		log.Errorf("Failed to Execute method : %s\n", se.step.Command.Method)
		return service.Result{Err: err}
	}
	log.Debugf("Executed method : %s.\n", se.step.Command.Method)

	if se.step.ResultFilter != "" {
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(res), &m); err != nil {
			log.Errorf("Failed json unmarshal : %v\n", err)
			return service.Result{Err: err}
		}

		o, err := jq.Process(m, se.step.ResultFilter)
		if err != nil {
			log.Errorf("Failed jq process : %v\n", err)
			return service.Result{Err: err}
		}

		b, err := json.Marshal(o)
		if err != nil {
			log.Errorf("Failed json marshal from jq process result : %v\n", err)
			return service.Result{Err: err}
		}
		return service.Result{Body: string(b)}
	}

	return service.Result{Body: res}
}
