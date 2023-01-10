package executor

import (
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
)

type ServiceExecutorV2 struct {
	service       service.ServiceV2
	updateChannel chan<- service.ServiceUpdateInterface
}

func NewServiceExecutorV2(service service.ServiceV2, updateChannel chan<- service.ServiceUpdateInterface) *ServiceExecutorV2 {
	return &ServiceExecutorV2{service: service, updateChannel: updateChannel}
}

func (se *ServiceExecutorV2) Execute() (err error) {
	var result service.Result
	flowStore := make(map[string]interface{})

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

	// update execute result to service scheduler through returnChannel.
	se.SendServiceStatusUpdate(0, service.StepStatusProcessing, "", se.service.StartTime, time.Time{})

	for i, step := range se.service.Flow {
		var te *StepExecutorV2

		if err := step.Inputs.FindReplacePassedInputsFrom(flowStore); err != nil {
			se.SendServiceStatusUpdate(i, service.StepStatusFail, err.Error(), se.service.StartTime, time.Now())
			return err
		}

		if stepInputs := step.Inputs.GetInputs(); stepInputs != nil {
			flowStore[step.Id+".inputs"] = step.Inputs.GetInputs()
		}

		te, err = NewStepExecutorV2(step)
		if err != nil {
			se.SendServiceStatusUpdate(i, service.StepStatusFail, err.Error(), se.service.StartTime, time.Now())
			return err
		}

		result = te.Execute()
		if err = result.Err; err != nil {
			se.SendServiceStatusUpdate(i, service.StepStatusFail, err.Error(), se.service.StartTime, time.Now())
			return err
		}
		flowStore[step.Id+".outputs"] = result.Body
	}

	// update execute result to service scheduler through returnChannel.
	se.SendServiceStatusUpdate(len(se.service.Flow)-1, service.StepStatusSuccess, result.Body, se.service.StartTime, time.Now())

	return nil
}

func (se *ServiceExecutorV2) SendServiceStatusUpdate(seq int, status service.StepStatus, result string, st, et time.Time) {
	if se.updateChannel != nil {
		update := service.UpdateServiceV2{
			Id:        se.service.Id,
			StepCount: len(se.service.Flow),
			Sequence:  seq,
			Status:    status,
			Result:    result,
			Started:   st,
			Ended:     et,
		}

		se.updateChannel <- &update
	}
}

type StepExecutorV2 struct {
	commander Commander
	step      *service.FlowStep
}

func NewStepExecutorV2(step *service.FlowStep) (*StepExecutorV2, error) {
	commander, err := NewCommander(&service.StepCommand{Method: step.Command, Args: step.Inputs.GetInputs()})
	if err != nil {
		return nil, err
	}

	return &StepExecutorV2{commander: commander, step: step}, nil
}

func (se *StepExecutorV2) Execute() service.Result {
	log.Debugf("Prepare to Execute method : %s\n", se.step.Command)
	res, err := se.commander.Run()
	if err != nil {
		log.Errorf("Failed to Execute method : %s\n", se.step.Command)
		return service.Result{Err: err}
	}
	log.Debugf("Executed method : %s.\n", se.step.Command)

	return service.Result{Body: res}
}
