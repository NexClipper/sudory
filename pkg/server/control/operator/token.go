package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
)

type Token struct {
	ctx database.Context
}

func NewToken(ctx database.Context) *Token {
	return &Token{ctx: ctx}
}

func (o *Token) Create(model tokenv1.Token) error {
	err := o.ctx.CreateToken(tokenv1.DbSchemaToken{Token: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Token) Get(uuid string) (*tokenv1.Token, error) {

	record, err := o.ctx.GetToken(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Token, nil
}

func (o *Token) Find(where string, args ...interface{}) ([]tokenv1.Token, error) {
	r, err := o.ctx.FindToken(where, args...)
	if err != nil {
		return nil, err
	}

	records := tokenv1.TransFormDbSchema(r)

	return records, nil
}

// func (o *Token) Query(cond *query_parser.QueryParser) ([]tokenv1.Token, error) {
// 	r, err := o.ctx.QueryToken(cond)
// 	if err != nil {
// 		return nil, err
// 	}

// 	records := tokenv1.TransFormDbSchema(r)

// 	return records, nil
// }

func (o *Token) Update(model tokenv1.Token) error {

	err := o.ctx.UpdateToken(tokenv1.DbSchemaToken{Token: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Token) Delete(uuid string) error {

	err := o.ctx.DeleteToken(uuid)
	if err != nil {
		return err
	}

	return nil
}
