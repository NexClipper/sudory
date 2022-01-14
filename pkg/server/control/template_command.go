package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/view"
	"github.com/labstack/echo/v4"
)

// Create Template Command
// @Description Create a template command
// @Accept json
// @Produce json
// @Tags server,server/template_command,create
// @Router /server/template_command [post]
// @Param template_command body v1.HttpReqTemplateCommand true "HttpReqTemplateCommand"
// @Success 200
func (c *Control) CreateTemplateCommand(ctx echo.Context) error {
	return view.NewCreateTemplateCommand(operator.NewCreateTemplateCommand(c.db)).Request(ctx)
}

// Get Template Command
// @Description Get a template command
// @Accept json
// @Produce json
// @Tags server,server/template_command,get
// @Router /server/template_command/{uuid} [get]
// @Param uuid path string true "HttpReqTemplateCommand 의 Uuid"
// @Success 200 {object} v1.HttpRspTemplate
func (c *Control) GetTemplateCommand(ctx echo.Context) error {
	return view.NewGetTemplateCommand(operator.NewGetTemplateCommand(c.db)).Request(ctx)
}

// Update Template Command
// @Description Update a template command
// @Accept json
// @Produce json
// @Tags server,server/template_command,update
// @Router /server/template_command [put]
// @Param template_command body v1.HttpReqTemplateCommand true "HttpReqTemplateCommand"
// @Success 200
func (c *Control) UpdateTemplateCommand(ctx echo.Context) error {
	return view.NewUpdateTemplateCommand(operator.NewUpdateTemplateCommand(c.db)).Request(ctx)
}

// Delete Template Command
// @Description Delete a template command
// @Accept json
// @Produce json
// @Tags server,server/template_command,delete
// @Router /server/template_command/{uuid} [delete]
// @Param uuid path string true "HttpReqTemplateCommand 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplateCommand(ctx echo.Context) error {
	return view.NewDeleteTemplateCommand(operator.NewDeleteTemplateCommand(c.db)).Request(ctx)
}
