package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

//Service
type Service struct {
	ctx database.Context
}

func NewService(ctx database.Context) *Service {
	return &Service{ctx: ctx}
}

func (o *Service) Create(model servicev1.Service) error {
	err := o.ctx.CreateService(servicev1.DbSchemaService{Service: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Service) Get(uuid string) (*servicev1.Service, error) {
	record, err := o.ctx.GetService(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Service, nil
}

func (o *Service) Find(where string, args ...interface{}) ([]servicev1.Service, error) {
	r, err := o.ctx.FindService(where, args...)
	if err != nil {
		return nil, err
	}

	records := servicev1.TransFormDbSchema(r)

	return records, nil
}

func (o *Service) Update(model servicev1.Service) error {
	err := o.ctx.UpdateService(servicev1.DbSchemaService{Service: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Service) Delete(uuid string) error {

	err := o.ctx.DeleteService(uuid)
	if err != nil {
		return err
	}

	//Service Step 레코드 삭제
	where := "service_uuid = ?"
	record, err := o.ctx.FindServiceStep(where, uuid)
	if err != nil {
		return err
	}

	for _, it := range record {
		err := o.ctx.DeleteServiceStep(it.Uuid)
		if err != nil {
			return err
		}
	}

	return nil
}

// type Service struct {
// 	db *database.DBManipulator

// 	ID        uint64
// 	Name      string
// 	ClusterID uint64
// 	StepCount uint
// 	Steps     []*Step

// 	Response ResponseFn
// }

// type Step struct {
// 	ID        uint64
// 	Name      string
// 	Sequence  uint64
// 	Command   string
// 	Parameter string
// }

// func NewService(d *database.DBManipulator) Operator {
// 	return &Service{db: d}
// }

// func (o *Service) toModelService() *model.Service {
// 	m := &model.Service{
// 		Name:      o.Name,
// 		ClusterID: o.ClusterID,
// 		StepCount: o.StepCount,
// 	}

// 	return m
// }

// func (o *Service) toModelStep(serviceID uint64) []*model.Step {
// 	var m []*model.Step
// 	for _, s := range o.Steps {
// 		modelStep := &model.Step{
// 			Name:      s.Name,
// 			ServiceID: serviceID,
// 			Sequence:  s.Sequence,
// 			Command:   s.Command,
// 			Parameter: s.Parameter,
// 		}

// 		m = append(m, modelStep)
// 	}

// 	return m
// }

// func (o *Service) Create(ctx echo.Context) error {
// 	service := o.toModelService()

// 	_, err := o.db.CreateService(service)
// 	if err != nil {
// 		return err
// 	}

// 	steps := o.toModelStep(service.ID)
// 	_, err = o.db.CreateStep(steps)
// 	if err != nil {
// 		return err
// 	}

// 	if o.Response != nil {
// 		o.Response(ctx, nil)
// 	}

// 	return nil
// }

// func (o *Service) Get(ctx echo.Context) error {
// 	service := o.toModelService()

// 	serviceSteps, err := o.db.GetServiceSteps(service)
// 	if err != nil {
// 		return err
// 	}

// 	if o.Response != nil {
// 		if len(serviceSteps) == 0 {
// 			o.Response(ctx, nil)
// 		} else {
// 			respBody := &model.RespService{
// 				Name:      serviceSteps[0].Service.Name,
// 				ClusterID: serviceSteps[0].Service.ClusterID,
// 				StepCount: serviceSteps[0].Service.StepCount,
// 			}
// 			for _, s := range serviceSteps {
// 				step := &model.RespStep{
// 					Name:      s.Step.Name,
// 					Sequence:  s.Step.Sequence,
// 					Command:   s.Step.Command,
// 					Parameter: s.Step.Parameter,
// 				}
// 				respBody.Step = append(respBody.Step, step)
// 			}
// 			o.Response(ctx, respBody)
// 		}
// 	}
// 	return nil
// }
