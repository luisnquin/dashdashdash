package echox

import (
	"github.com/labstack/echo/v4"
)

type ControllersGetter interface {
	GetControllers() []Controller
}

type Controller struct {
	// The HTTP method name to access to this controller.
	Method string
	// The underlying request handler.
	Handler echo.HandlerFunc
	// The path that will correspond to the controller.
	Path string
	// The middlewares to be applied after the auth+verbose middleware.
	Middlewares []echo.MiddlewareFunc
	// Add logs when a request is reached by the controller.
	Verbose bool
	// Add Authentication based on JWT + OTP on the current controller.
	Auth bool
}

func LoadControllers(e *echo.Echo, authMiddleware echo.MiddlewareFunc, getters []ControllersGetter) {
	for _, item := range getters {
		for _, controller := range item.GetControllers() {
			if controller.Auth {
				controller.Middlewares = append([]echo.MiddlewareFunc{authMiddleware}, controller.Middlewares...)
			}

			if controller.Verbose {
				controller.Middlewares = append([]echo.MiddlewareFunc{PreAuthVerboseRequestMiddleware}, controller.Middlewares...)
				controller.Middlewares = append(controller.Middlewares, PostAuthVerboseRequestMiddleware)
			}

			controller.Middlewares = append([]echo.MiddlewareFunc{FirstMiddleware}, controller.Middlewares...)

			e.Add(controller.Method, controller.Path, controller.Handler, controller.Middlewares...)
		}
	}
}
