package control

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/block"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Control struct {
	db *database.DBManipulator
}

func New(d *database.DBManipulator) *Control {
	return &Control{db: d}
}

func (ctl Control) Scope(fn func(database.Context) (interface{}, error)) (v interface{}, err error) {
	block.Block{
		Try: func() {
			_, lockerr := ctl.db.Engine().Transaction(func(s *xorm.Session) (interface{}, error) {
				v, err = fn(database.NewXormContext(s))
				return nil, err
			})
			if err == nil && lockerr != nil {
				err = errors.Wrapf(lockerr, "xorm commit")
			}
		},
		Catch: func(ex error) {
			err = errors.Wrapf(ex, "catch")
		},
	}.Do()

	return
}

func (ctl Control) NewSession() database.Context {
	return database.NewXormContext(ctl.db.Engine().NewSession())
}
