package systemd

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

type (
	Module struct {
		repo moduleRepository
	}

	moduleRepository struct {
		systemd Repository
	}
)

func NewModule() Module {
	return Module{
		repo: moduleRepository{
			systemd: NewRepository(),
		},
	}
}

func (m Module) GetControllers() []echox.Controller {
	return []echox.Controller{
		{
			Method:  http.MethodGet,
			Path:    "/host/systemd/units",
			Handler: m.ListUnitsHandler(),
			Middlewares: []echo.MiddlewareFunc{
				m.IsRunningMiddleware(),
			},
			Auth: true,
		},
	}
}
