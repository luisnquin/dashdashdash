package echox

import (
	"github.com/labstack/echo/v4"
)

type ControllersGetter interface {
	GetControllers() []Controller
}

func LoadControllers(e *echo.Echo, getters []ControllersGetter) {
	for _, item := range getters {
		// e.Group() <- global middlewares

		for _, controller := range item.GetControllers() {
			controller.Materialize(e)
		}
	}
}

type Controller struct {
	Method      string
	Handler     echo.HandlerFunc
	Path        string
	Middlewares []echo.MiddlewareFunc
}

func (c Controller) Materialize(e *echo.Echo) {
	e.Add(c.Method, c.Path, c.Handler, c.Middlewares...)
}
