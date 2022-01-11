package database

import "github.com/NexClipper/sudory/pkg/server/model"

func (d *DBManipulator) CreateToken(m *model.Token) (int64, error) {
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
