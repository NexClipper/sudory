package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/view"
	"github.com/labstack/echo/v4"
)

// GetCatalogue
// @Description Get a Catalogues
// @Produce json
// @Tags server
// @Router /server/catalogue [get]
// @Success 200 {object} model.Catalogues
func (c *Control) GetCatalogue(ctx echo.Context) error {
	v := view.NewGetCatalogue(operator.NewCatalogue())
	return v.Request(ctx)
}
