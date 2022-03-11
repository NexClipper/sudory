package vault

// import (
// 	"github.com/NexClipper/sudory/pkg/server/model"
// 	"github.com/labstack/echo/v4"
// )

// type Catalogue struct {
// 	Name string

// 	Response ResponseFn
// }

// func NewCatalogue() Operator {
// 	return &Catalogue{}
// }

// func (o *Catalogue) Create(ctx echo.Context) error {
// 	return nil
// }

// func (o *Catalogue) Get(ctx echo.Context) error {
// 	m := &model.Catalogues{}

// 	c := &model.Catalogue{Name: "GET_NAMESPACE"}
// 	m.Items = append(m.Items, c)

// 	if o.Response != nil {
// 		return o.Response(ctx, m)
// 	}

// 	return nil
// }
