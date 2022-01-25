package service

import (
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
)

type ServiceExecutor struct {
	service       Service
	updateChannel chan<- Service
}

func NewServiceExecutor(service Service, updateChannel chan<- Service) *ServiceExecutor {
	return &ServiceExecutor{service: service, updateChannel: updateChannel}
}

func (se *ServiceExecutor) Execute() (err error) {
	var result Result

	defer func() {
		if err != nil {
			log.Errorf(err.Error())
			se.service.Status = ServiceStatusFailed
			se.service.Result = result
			se.SendStatusUpdate()
		} else {
			se.service.Status = ServiceStatusSuccess
			se.service.Result = result
			se.SendStatusUpdate()
		}
	}()

	se.service.Status = ServiceStatusProcessing
	se.SendStatusUpdate()
	for i, step := range se.service.Steps {
		var te *StepExecutor

		te, err = NewStepExecutor(*step)
		if err != nil {
			return err
		}
		// update execute result to service scheduler through returnChannel.
		se.service.Steps[i].Status = StepStatusProcessing
		se.service.Steps[i].StartTime = time.Now()
		se.SendStatusUpdate()
		result = te.Execute()
		if err = result.err; err != nil {
			se.service.Steps[i].Status = StepStatusFail
			se.service.Steps[i].EndTime = time.Now()
			return err
		}

		// update execute result to service scheduler through returnChannel.
		se.service.Steps[i].Status = StepStatusSuccess
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
	step      Step
}

func NewStepExecutor(step Step) (*StepExecutor, error) {
	commander, err := NewCommander(step.Command)
	if err != nil {
		return nil, err
	}

	return &StepExecutor{commander: commander, step: step}, nil
}

func (se *StepExecutor) Execute() Result {
	res, err := se.commander.Run()
	if err != nil {
		log.Errorf("Failed to Execute method : %s.\n", se.step.Command.Method)
		return Result{err: err}
	}
	log.Debugf("Executed method : %s.\n", se.step.Command.Method)

	return Result{body: res}
}
