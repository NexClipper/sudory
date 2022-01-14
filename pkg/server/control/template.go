package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/view"
	"github.com/labstack/echo/v4"
)

// Create Template
// @Description Create []template
// @Accept json
// @Produce json
// @Tags server,server/template,create
// @Router /server/template [post]
// @Param template body v1.HttpReqTemplates true "HttpReqTemplates"
// @Success 200
func (c *Control) CreateTemplate(ctx echo.Context) error {
	return view.NewCreateTemplate(operator.NewCreateTemplate(c.db)).Request(ctx)
}

// Get Template
// @Description Get a template
// @Accept json
// @Produce json
// @Tags server,server/template,get
// @Router /server/template/{uuid} [get]
// @Param uuid path string true "ReqTemplate 의 Uuid"
// @Success 200 {object} v1.HttpRspTemplate
func (c *Control) GetTemplate(ctx echo.Context) error {
	return view.NewGetTemplate(operator.NewGetTemplate(c.db)).Request(ctx)
}

// Find []Template
// @Description Find []template
// @Accept json
// @Produce json
// @Tags server,server/template,find
// @Router /server/template [get]
// @Param uuid   query string false "ReqTemplate 의 Uuid"
// @Param name   query string false "ReqTemplate 의 Name"
// @Param origin query string false "ReqTemplate 의 Origin"
// @Success 200 {object} []v1.HttpRspTemplate "[]HttpRspTemplate"
func (c *Control) FindTemplate(ctx echo.Context) error {
	return view.NewFindTemplate(operator.NewFindTemplate(c.db)).Request(ctx)
}

// Update Template
// @Description Update a template
// @Accept json
// @Produce json
// @Tags server,server/template,update
// @Router /server/template [put]
// @Param template body v1.HttpReqTemplate true "HttpReqTemplate"
// @Success 200
func (c *Control) UpdateTemplate(ctx echo.Context) error {
	return view.NewUpdateTemplate(operator.NewUpdateTemplate(c.db)).Request(ctx)
}

// Delete Template
// @Description Delete a template
// @Accept json
// @Produce json
// @Tags server,server/template,delete
// @Router /server/template/{uuid} [delete]
// @Param uuid path string true "HttpReqTemplate 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplate(ctx echo.Context) error {
	return view.NewDeleteTemplate(operator.NewDeleteTemplate(c.db)).Request(ctx)
}
