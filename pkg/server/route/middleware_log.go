package route

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/labstack/echo/v4"
)

func echoServiceLogger(w io.Writer) echo.MiddlewareFunc {
	//echo logger
	format := fmt.Sprintf("{%v}\n",
		strings.Join([]string{
			`"time":"${time_rfc3339_nano}"`,
			`"id":"${id}"`,
			`"remote_ip":"${remote_ip}"`,
			`"host":"${host}"`,
			`"method":"${method}"`,
			`"uri":"${uri}"`,
			`"status":${status}`,
			`"error":"${error}"`,
			`"latency":${latency}`,
			`"latency_human":"${latency_human}"`,
			`"bytes_in":${bytes_in}`,
			`"bytes_out":${bytes_out}`,
		}, ","))

	logconfig := DefaultLoggerConfig
	logconfig.Output = w
	logconfig.Format = format

	return LoggerWithConfig(logconfig)
}

func echoErrorResponder(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if httperr, ok := err.(*echo.HTTPError); ok {
		code = httperr.Code
		if httperr.Internal != nil {
			err = httperr.Internal
		}
	}

	ctx.JSON(code, map[string]interface{}{
		"code": code,
		// "status":     http.StatusText(code),
		"message": err.Error(),
	})
}

func echoErrorLogger(err error, ctx echo.Context) {
	nullstring := func(p *string) (s string) {
		s = fmt.Sprintf("%v", p)
		if p != nil {
			s = *p
		}
		return
	}

	code := http.StatusInternalServerError
	if httperr, ok := err.(*echo.HTTPError); ok {
		code = httperr.Code
		if httperr.Internal != nil {
			err = httperr.Internal
		}
	}

	var stack *string
	//stack for surface
	logs.StackIter(err, func(s string) {
		stack = &s
	})
	//stack for internal
	logs.CauseIter(err, func(err error) {
		logs.StackIter(err, func(s string) {
			stack = &s
		})
	})

	id := ctx.Response().Header().Get(echo.HeaderXRequestID)

	reqbody, _ := echoutil.Body(ctx)

	logger.Error(fmt.Errorf("%v%v", err.Error(), logs.KVL(
		"id", id,
		"code", code,
		"reqbody", reqbody,
		"stack", nullstring(stack),
	)))
}
