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
	updateChannel chan<- service.Service
}

func NewServiceExecutor(service service.Service, updateChannel chan<- service.Service) *ServiceExecutor {
	return &ServiceExecutor{service: service, updateChannel: updateChannel}
}

func (se *ServiceExecutor) Execute() (err error) {
	var result service.Result

	defer func() {
		se.service.EndTime = time.Now()
		if err != nil {
			log.Errorf(err.Error())
			se.service.Status = service.ServiceStatusFailed
			se.service.Result.Err = err
			se.SendStatusUpdate()
		} else {
			se.service.Status = service.ServiceStatusSuccess
			se.service.Result = result
			se.SendStatusUpdate()
		}
	}()

	se.service.StartTime = time.Now()
	se.service.Status = service.ServiceStatusProcessing
	se.SendStatusUpdate()
	for i, step := range se.service.Steps {
		var te *StepExecutor

		te, err = NewStepExecutor(*step)
		if err != nil {
			se.service.Steps[i].Status = service.StepStatusFail
			se.service.Steps[i].EndTime = time.Now()
			return err
		}
		// update execute result to service scheduler through returnChannel.
		se.service.Steps[i].Status = service.StepStatusProcessing
		se.service.Steps[i].StartTime = time.Now()
		se.SendStatusUpdate()
		result = te.Execute()
		if err = result.Err; err != nil {
			se.service.Steps[i].Status = service.StepStatusFail
			se.service.Steps[i].EndTime = time.Now()
			return err
		}

		// update execute result to service scheduler through returnChannel.
		se.service.Steps[i].Status = service.StepStatusSuccess
		se.service.Steps[i].EndTime = time.Now()
		se.SendStatusUpdate()
	}

	return nil
}

func (se *ServiceExecutor) SendStatusUpdate() {
	if se.updateChannel != nil {
		se.updateChannel <- se.service
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
	res, err := se.commander.Run()
	if err != nil {
		log.Errorf("Failed to Execute method : %s.\n", se.step.Command.Method)
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
