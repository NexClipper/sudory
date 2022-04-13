package control

import (
	"bytes"
	"net/http"
	"sort"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// FindTemplateRecipe
// @Description Find TemplateRecipe
// @Produce     json
// @Tags        server/template_recipe
// @Router      /server/template_recipe [get]
// @Param       x_auth_token header string false "client session token"
// @Param       method       query  string false "Template Command 의 Method"
// @Success     200 {array} v1.TemplateRecipe
func (ctl Control) FindTemplateRecipe(ctx echo.Context) error {
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

	cond :=
		prepare.WrapMap("like",
			prepare.WrapMap("method", buff.String()))

	records, err := vault.NewTemplateRecipe(ctl.NewSession()).Prepare(cond)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find template recipe"))
	}

	//sort by args
	sort.Slice(records, func(i, j int) bool {
		return nullable.String(records[i].Args).Value() < nullable.String(records[j].Args).Value()
	})
	//sort by name
	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})

	return ctx.JSON(http.StatusOK, records)
}
