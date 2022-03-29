package control

import (
	"bytes"
	"sort"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	recipev1 "github.com/NexClipper/sudory/pkg/server/model/template_recipe/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// FindTemplateRecipe
// @Description Find TemplateRecipe
// @Produce     json
// @Tags        server/template_recipe
// @Router      /server/template_recipe [get]
// @Param       method query string false "Template Command 의 Method"
// @Success     200 {array} v1.HttpRspTemplateRecipe
func (c *Control) FindTemplateRecipe() func(ctx echo.Context) error {
	binder := func(ctx Context) error {
		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		method := ctx.Queries()["method"]

		buff := bytes.Buffer{}
		for i, s := range strings.Split(method, ".") {
			if 0 < i {
				buff.WriteString(".")
			}
			buff.WriteString(s)
		}
		//뒤에 like 조회 와일드 카드를 붙여준다
		buff.WriteString("%")

		cond :=
			prepare.WrapMap("like",
				prepare.WrapMap("method", buff.String()))

		records, err := vault.NewTemplateRecipe(ctx.Database()).Prepare(cond)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateRecipe Query")
		}

		//sort by args
		sort.Slice(records, func(i, j int) bool {
			return nullable.String(records[i].Args).Value() < nullable.String(records[j].Args).Value()
		})
		//sort by name
		sort.Slice(records, func(i, j int) bool {
			return nullable.String(records[i].Name).Value() < nullable.String(records[j].Name).Value()
		})

		return recipev1.TransToHttpRsp(records), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "FindTemplateRecipe binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindTemplateRecipe operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}
