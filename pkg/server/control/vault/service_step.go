package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/pkg/errors"
)

//ServiceStep
type ServiceStep struct {
	ctx database.Context
}

func NewServiceStep(ctx database.Context) *ServiceStep {
	return &ServiceStep{ctx: ctx}
}

func (vault ServiceStep) Create(model stepv1.ServiceStep) (*stepv1.ServiceStep, error) {
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
			))
	}

	return &model, nil
}

func (vault ServiceStep) Get(uuid string) (*stepv1.ServiceStep, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &stepv1.ServiceStep{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
				"where", where,
				"args", args,
			))
	}

	return model, nil
}

func (vault ServiceStep) Find(where string, args ...interface{}) ([]stepv1.ServiceStep, error) {
	models := make([]stepv1.ServiceStep, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
				"where", where,
				"args", args,
			))
	}

	return models, nil
}

func (vault ServiceStep) Query(query map[string]string) ([]stepv1.ServiceStep, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]stepv1.ServiceStep, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
				"query", query,
			))
	}

	return models, nil
}

func (vault ServiceStep) Update(model stepv1.ServiceStep) (*stepv1.ServiceStep, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	if err := vault.ctx.Where(where, args...).Update(&model); err != nil {
		return nil, errors.Wrapf(err, "database update%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
				"where", where,
				"args", args,
			))
	}

	return &model, nil
}

func (vault ServiceStep) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &stepv1.ServiceStep{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
				"where", where,
				"args", args,
			))
	}

	return nil
}

func (vault ServiceStep) Delete_ByService(service_uuid string) error {
	where := "service_uuid = ?"
	args := []interface{}{
		service_uuid,
	}
	model := &stepv1.ServiceStep{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"table", new(stepv1.ServiceStep).TableName(),
				"where", where,
				"args", args,
			))
	}

	return nil
}

// // ChainingSequence
// //  uuid: 해당 객체는 대상에서 제외
// //  대상 객체 외는 순서에 맞추어 Sequence 지정
// func (vault ServiceStep) ChainingSequence(service_uuid, uuid string) error {
// 	where := "service_uuid = ?"
// 	steps, err := vault.Find(where, service_uuid)
// 	if err != nil {
// 		return errors.Wrapf(err, "Database Find")
// 	}

// 	//sort -> Sequence
// 	sort.Slice(steps, func(i, j int) bool {
// 		return nullable.Int32(steps[i].Sequence).Value() < nullable.Int32(steps[j].Sequence).Value()
// 	})

// 	seq := int32(0)
// 	for i := range steps {
// 		if steps[i].Uuid != uuid {
// 			steps[i].Sequence = newist.Int32(int32(seq))
// 		}
// 		seq++
// 	}
// 	for i := range steps {
// 		if _, err := vault.Update(steps[i]); err != nil {
// 			return errors.Wrapf(err, "Database Update")
// 		}
// 	}

// 	return nil
// }
