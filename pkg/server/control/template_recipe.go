package control

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	recipev2 "github.com/NexClipper/sudory/pkg/server/model/template_recipe/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Find TemplateRecipe
// @Security    XAuthToken
// @Produce     json
// @Tags        server/template_recipe
// @Router      /server/template_recipe [get]
// @Param       method       query  string false "Template Command Method"
// @Success     200 {array} v2.HttpRsp_TemplateRecipe
func (ctl ControlVanilla) FindTemplateRecipe(ctx echo.Context) (err error) {
	method := echoutil.QueryParam(ctx)["method"]
	buff := bytes.Buffer{}
	for i, s := range strings.Split(method, ".") {
		if 0 < i {
			buff.WriteString(".")
		}
		buff.WriteString(s)
	}
	//뒤에 like 조회 와일드 카드를 붙여준다
	buff.WriteString("%")

	var p = stmt.Limit(50, 1)
	if 0 < len(echoutil.QueryParam(ctx)["p"]) {
		p, err = stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
		err = errors.Wrapf(err, "failed to parse pagination")
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	}

	rsp := make([]recipev2.HttpRsp_TemplateRecipe, 0, state.ENV__INIT_SLICE_CAPACITY__())

	recipe := recipev2.TemplateRecipe{}
	recipe.Method = buff.String()
	like_method := stmt.Like("method", recipe.Method)
	order := stmt.Asc("name", "args")

	err = stmtex.Select(recipe.TableName(), recipe.ColumnNames(), like_method, order, p).
		QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner, _ int) (err error) {
			err = recipe.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "failed to scan")
			}
			rsp = append(rsp, recipe)
			return
		})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, []recipev2.HttpRsp_TemplateRecipe(rsp))
}
