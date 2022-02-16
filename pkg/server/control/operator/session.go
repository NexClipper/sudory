package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

type Session struct {
	ctx database.Context
}

func NewSession(ctx database.Context) *Session {
	return &Session{ctx: ctx}
}

func (o *Session) Create(model sessionv1.Session) error {
	err := o.ctx.CreateSession(sessionv1.DbSchemaSession{Session: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Session) Get(uuid string) (*sessionv1.Session, error) {

	record, err := o.ctx.GetSession(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Session, nil
}

func (o *Session) Find(where string, args ...interface{}) ([]sessionv1.Session, error) {
	r, err := o.ctx.FindSession(where, args)
	if err != nil {
		return nil, err
	}

	records := sessionv1.TransFormDbSchema(r)

	return records, nil
}

// func (o *Session) Query(cond *query_parser.QueryParser) ([]sessionv1.Session, error) {
// 	r, err := o.ctx.QuerySession(cond)
// 	if err != nil {
// 		return nil, err
// 	}

// 	records := sessionv1.TransFormDbSchema(r)

// 	return records, nil
// }

func (o *Session) Update(model sessionv1.Session) error {

	err := o.ctx.UpdateSession(sessionv1.DbSchemaSession{Session: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Session) Delete(uuid string) error {

	err := o.ctx.DeleteSession(uuid)
	if err != nil {
		return err
	}

	return nil
}
