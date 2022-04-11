package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	"github.com/pkg/errors"
)

//Service
type Service struct {
	ctx database.Context
}

func NewService(ctx database.Context) *Service {
	return &Service{ctx: ctx}
}

func (vault Service) Create(model servicev1.Service) (*servicev1.Service, error) {
	//create service
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return &model, nil
}

func (vault Service) Get(uuid string) (*servicev1.Service, error) {
	//get service
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &servicev1.Service{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	// //find step
	// where = "service_uuid = ?"
	// args = []interface{}{
	// 	service.Uuid,
	// }
	// steps := make([]stepv1.ServiceStep, 0)
	// if err := vault.ctx.Where(where, args...).Find(&steps); err != nil {
	// 	return nil, errors.Wrapf(err, "database find%v",
	// 		logs.KVL(
	// 			"where", where,
	// 			"args", args,
	// 		))
	// }
	// //sort -> Sequence ASC
	// sort.Slice(steps, func(i, j int) bool {
	// 	return nullable.Int32(steps[i].Sequence).Value() < nullable.Int32(steps[j].Sequence).Value()
	// })

	return model, nil
}

func (vault Service) Find(where string, args ...interface{}) ([]servicev1.Service, error) {
	//find service
	models := make([]servicev1.Service, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}
	// //make result
	// var services = make([]servicev1.DbSchemaServiceAndSteps, len(records))
	// for i := range records {
	// 	service := records[i]
	// 	//set service
	// 	services[i].Service = service
	// 	//find step
	// 	where := "service_uuid = ?"
	// 	args := []interface{}{
	// 		service.Uuid,
	// 	}
	// 	steps := make([]stepv1.ServiceStep, 0)
	// 	if err := vault.ctx.Where(where, args...).Find(&steps); err != nil {
	// 		return nil, errors.Wrapf(err, "database find%v",
	// 			logs.KVL(
	// 				"where", where,
	// 				"args", args,
	// 			))
	// 	}
	// 	//sort -> Sequence ASC
	// 	sort.Slice(steps, func(i, j int) bool {
	// 		return nullable.Int32(steps[i].Sequence).Value() < nullable.Int32(steps[j].Sequence).Value()
	// 	})

	// 	//set steps
	// 	services[i].Steps = steps
	// }

	return models, nil
}

func (vault Service) Query(query map[string]string) ([]servicev1.Service, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]servicev1.Service, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	// //make result
	// var services = make([]servicev1.DbSchemaServiceAndSteps, len(records))
	// for i := range records {
	// 	service := records[i]
	// 	//set service
	// 	services[i].Service = service
	// 	//find step
	// 	where := "service_uuid = ?"
	// 	args := []interface{}{
	// 		service.Uuid,
	// 	}
	// 	steps := make([]stepv1.ServiceStep, 0)
	// 	if err := vault.ctx.Where(where, args...).Find(&steps); err != nil {
	// 		return nil, errors.Wrapf(err, "database find%v",
	// 			logs.KVL(
	// 				"where", where,
	// 				"args", args,
	// 			))
	// 	}
	// 	//sort -> Sequence ASC
	// 	sort.Slice(steps, func(i, j int) bool {
	// 		return nullable.Int32(steps[i].Sequence).Value() < nullable.Int32(steps[j].Sequence).Value()
	// 	})

	// 	//set steps
	// 	services[i].Steps = steps
	// }

	return models, nil
}

func (vault Service) Update(model servicev1.Service) (*servicev1.Service, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	// record := &servicev1.Service{Service: model}
	if err := vault.ctx.Where(where, args...).Update(&model); err != nil {
		return nil, errors.Wrapf(err, "database update%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return &model, nil
}

func (vault Service) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &servicev1.Service{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}
	return nil
}
