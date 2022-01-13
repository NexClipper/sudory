package operator

import (
	"github.com/NexClipper/sudory/pkg/server/model"
	"github.com/labstack/echo/v4"
)

type ResponseFn func(ctx echo.Context, m model.Modeler) error

type Operator interface {
	Create(ctx echo.Context) error
	Get(ctx echo.Context) error
}

//생성
type Creator interface {
	Create(ctx echo.Context) error
}

//조회
type Getter interface {
	Get(ctx echo.Context) error
}

//갱신
type Updater interface {
	Update(ctx echo.Context) error
}

//삭제
type Remover interface {
	Delete(ctx echo.Context) error
}
