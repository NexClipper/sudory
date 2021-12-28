package control

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/control/operator"
	"github.com/NexClipper/sudory-prototype-r1/pkg/view"
	"github.com/labstack/echo/v4"
)

// GetCatalogue
// @Description Get a Catalogues
// @Produce json
// @Router /catalogue [get]
// @Success 200 {object} model.Catalogues
func (c *Control) GetCatalogue(ctx echo.Context) error {
	v := view.NewGetCatalogue(operator.NewCatalogue())
	return v.Request(ctx)
}
