package database

import "github.com/NexClipper/sudory-prototype-r1/pkg/model"

func (d *DBManipulator) CreateService(m *model.Service) (int64, error) {
	tx := d.session()
	tx.Begin()

	cnt, err := tx.Insert(m)

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return cnt, err
}

func (d *DBManipulator) GetService(m *model.Service) (*model.Service, error) {
	tx := d.session()

	var service model.Service
	_, err := tx.ID(m.ID).Get(&service)

	return &service, err
}

func (d *DBManipulator) CreateStep(m []*model.Step) (int64, error) {
	tx := d.session()
	tx.Begin()

	cnt, err := tx.InsertMulti(m)

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return cnt, err
}
