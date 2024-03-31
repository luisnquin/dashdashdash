package systemd

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

type Module struct{}

func NewModule() Module {
	return Module{}
}

func (m Module) GetControllers() []echox.Controller {
	return []echox.Controller{
		{
			Method:  http.MethodGet,
			Path:    "/host/systemd/services",
			Handler: m.ListServicesHandler(),
			Middlewares: []echo.MiddlewareFunc{
				m.IsRunningMiddleware(),
			},
		},
	}
}

// func (m Module) GetGlobalMiddlewares() []echo.MiddlewareFunc {
// 	return []echo.MiddlewareFunc{
// 		m.IsRunningMiddleware(),
// 	}
// }
