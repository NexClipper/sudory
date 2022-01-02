package operator

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/database"
	"github.com/NexClipper/sudory-prototype-r1/pkg/model"
	"github.com/labstack/echo/v4"
)

type Token struct {
	db *database.DBManipulator

	ClusterID uint64
	Key       string

	Response ResponseFn
}

func NewToken(d *database.DBManipulator) Operator {
	return &Token{db: d}
}

func (o *Token) toModel() *model.Token {
	m := &model.Token{
		ClusterID: o.ClusterID,
		Key:       o.Key,
	}

	return m
}

func (o *Token) Create(ctx echo.Context) error {
	token := o.toModel()

	_, err := o.db.CreateToken(token)
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

func (o *Token) Get(ctx echo.Context) error {
	return nil
}
