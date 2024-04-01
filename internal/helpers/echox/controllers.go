package echox

import (
	"github.com/labstack/echo/v4"
)

type ControllersGetter interface {
	GetControllers() []Controller
}

type Controller struct {
	Method      string
	Handler     echo.HandlerFunc
	Path        string
	Middlewares []echo.MiddlewareFunc
	Auth        bool
}

func LoadControllers(e *echo.Echo, authMiddleware echo.MiddlewareFunc, getters []ControllersGetter) {
	for _, item := range getters {
		for _, controller := range item.GetControllers() {
			if controller.Auth {
				controller.Middlewares = append([]echo.MiddlewareFunc{authMiddleware}, controller.Middlewares...)
			}

			e.Add(controller.Method, controller.Path, controller.Handler, controller.Middlewares...)
		}
	}
}
