package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type ServiceStep struct {
	tx *xorm.Session
}

func NewServiceStep(tx *xorm.Session) *ServiceStep {
	return &ServiceStep{tx: tx}
}

// Create
func (vault ServiceStep) Create(record stepv1.ServiceStep) (*stepv1.ServiceStep, error) {
	if err := database.XormCreate(
		vault.tx, &record); err != nil {
		return nil, errors.Wrapf(err, "create %v", record.TableName())
	}

	return &record, nil
}

// Find
func (vault ServiceStep) Find(where string, args ...interface{}) ([]stepv1.ServiceStep, error) {
	record := make([]stepv1.ServiceStep, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &record); err != nil {
		return record, errors.Wrapf(err, "find %v", new(stepv1.ServiceStep).TableName())
	}

	return record, nil
}

// Get
func (vault ServiceStep) Get(uuid string) (*stepv1.ServiceStep, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &stepv1.ServiceStep{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), record); err != nil {
		return record, errors.Wrapf(err, "get %v", new(stepv1.ServiceStep).TableName())
	}

	return record, nil
}

// Update
func (vault ServiceStep) Update(record stepv1.ServiceStep) (*stepv1.ServiceStep, error) {
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
