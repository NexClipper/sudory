package echoutil

import "github.com/labstack/echo/v4"

// HttpError
func HttpError(err error, code int) error {
	if err == nil {
		return nil
	}
	return echo.NewHTTPError(code).SetInternal(err)
}

func WrapHttpError(code int, fn ...func() error) error {
	for _, fn := range fn {
		if err := fn(); err != nil {
			return echo.NewHTTPError(code).SetInternal(err)
		}
	}
	return nil
}
