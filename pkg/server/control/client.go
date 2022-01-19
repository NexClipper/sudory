package control

// import (
// 	"github.com/NexClipper/sudory/pkg/server/control/operator"
// 	"github.com/NexClipper/sudory/pkg/server/view"
// 	"github.com/labstack/echo/v4"
// )

// // CreateClient
// // @Description Regist a Client
// // @Accept json
// // @Produce json
// // @Tags client
// // @Router /client/regist [post]
// // @Param client body model.ReqClient true "Client의 정보"
// // @Success 200
// func (c *Control) CreateClient(ctx echo.Context) error {
// 	v := view.NewCreateClient(operator.NewClient(c.db))
// 	return v.Request(ctx)
// }
