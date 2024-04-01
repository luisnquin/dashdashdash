package echox

import (
	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/luisnquin/go-log"
)

// Request ID placeholder.
const REQUEST_ID = "request_id"

func FirstMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		g, err := nanoid.Standard(12)
		if err != nil {
			panic(err)
		}

		c.Set(REQUEST_ID, g())

		return next(c)
	}
}

func PreAuthVerboseRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		rid, _ := c.Get(REQUEST_ID).(string)

		log.Info().Str("rid", rid).Str("method", r.Method).Str("path", r.URL.Path).Str("agent", r.UserAgent()).Send()

		return next(c)
	}
}

func PostAuthVerboseRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if ok {
			rid, _ := c.Get(REQUEST_ID).(string)

			log.Info().Str("rid", rid).Any("user", map[string]string{"email": user.Email, "password": user.Password}).Send()
		}

		return next(c)
	}
}
