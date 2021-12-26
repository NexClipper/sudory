package operator

import "github.com/labstack/echo/v4"

type Cluster struct {
	Name string

	Response ResponseFn
}

func NewCluster() Operator {
	return &Cluster{}
}

func (o *Cluster) Create(ctx echo.Context) error {

	return nil
}
