package vault

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Service struct {
	tx *xorm.Session
}

func NewService(tx *xorm.Session) *Service {
	return &Service{tx: tx}
}

func (vault Service) Create(record servicev1.Service) (*servicev1.Service, error) {
	//create service
	if err := database.XormCreate(
		vault.tx, &record); err != nil {
		return nil, errors.Wrapf(err, "create service")
	}

	return &record, nil
}

func (vault Service) Get(uuid string) (*servicev1.Service, []stepv1.ServiceStep, error) {

	//get service
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	service := &servicev1.Service{}
	if err := database.XormGet(vault.tx.Where(where, args...), service); err != nil {
		return nil, nil, errors.Wrapf(err, "get service")
	}

	//find step
	steps, err := NewServiceStep(vault.tx).Find("service_uuid = ?", service.Uuid)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "get service")
	}

	sort.Slice(steps, func(i, j int) bool {
		var a, b int32 = 0, 0
		if steps[i].Sequence != nil {
			a = *steps[i].Sequence
		}
		if steps[j].Sequence != nil {
			b = *steps[j].Sequence
		}
		return a < b
	})

	return service, steps, nil
}

func (vault Service) Find(where string, args ...interface{}) ([]servicev1.Service, map[string][]stepv1.ServiceStep, error) {
	services := make([]servicev1.Service, 0)
	stepSet := map[string][]stepv1.ServiceStep{}

	//find service
	if err := database.XormFind(
		vault.tx.Where(where, args...), &services); err != nil {
		return services, stepSet, errors.Wrapf(err, "find %v", new(servicev1.Service).TableName())
	}

	for _, service := range services {
		where := "service_uuid = ?"
		steps, err := NewServiceStep(vault.tx).Find(where, service.Uuid)
		if err != nil {
			return services, stepSet, errors.Wrapf(err, "find %v", new(servicev1.Service).TableName())
		}

		sort.Slice(steps, func(i, j int) bool {
			var a, b int32 = 0, 0
			if steps[i].Sequence != nil {
				a = *steps[i].Sequence
			}
			if steps[j].Sequence != nil {
				b = *steps[j].Sequence
			}
			return a < b
		})

		stepSet[service.Uuid] = steps
	}

	return services, stepSet, nil
}

func (vault Service) Query(query map[string]string) ([]servicev1.Service, map[string][]stepv1.ServiceStep, error) {

	services := make([]servicev1.Service, 0)
	stepSet := map[string][]stepv1.ServiceStep{}

	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return services, stepSet, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &services); err != nil {
		return services, stepSet, errors.Wrapf(err, "query %v", new(servicev1.Service).TableName())
	}

	for _, service := range services {
		where := "service_uuid = ?"
		steps, err := NewServiceStep(vault.tx).Find(where, service.Uuid)
		if err != nil {
			return services, stepSet, errors.Wrapf(err, "query %v", new(servicev1.Service).TableName())
		}

		sort.Slice(steps, func(i, j int) bool {
			var a, b int32 = 0, 0
			if steps[i].Sequence != nil {
				a = *steps[i].Sequence
			}
			if steps[j].Sequence != nil {
				b = *steps[j].Sequence
			}
			return a < b
		})

		stepSet[service.Uuid] = steps
	}

	return services, stepSet, nil
}

// Update
func (vault Service) Update(record servicev1.Service) (*servicev1.Service, error) {
	where := "uuid = ?"
	args := []interface{}{
		record.Uuid,
	}
	if err := database.XormUpdate(
		vault.tx.Where(where, args...), &record); err != nil {
		return nil, errors.Wrapf(err, "update %v", record.TableName())
	}

	return &record, nil
}

// Delete
func (vault Service) Delete(uuid string) error {
	//delete service steps
	step := new(stepv1.ServiceStep)
	if err := database.XormDelete(
		vault.tx.Where("service_uuid = ?", uuid), step); err != nil {
		return errors.Wrapf(err, "delete %v", step.TableName())
	}
	//delete service
	service := new(servicev1.Service)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), service); err != nil {
		return errors.Wrapf(err, "delete %v", service.TableName())
	}

	return nil
}
