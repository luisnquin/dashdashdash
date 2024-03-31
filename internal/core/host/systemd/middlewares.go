package systemd

import (
	"net/http"

	systemd_utils "github.com/coreos/go-systemd/v22/util"
	"github.com/labstack/echo/v4"
)

func (m Module) IsRunningMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !systemd_utils.IsRunningSystemd() {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "systemd is not running in current host")
			}

			return next(c)
		}
	}
}
