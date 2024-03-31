package core

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/core/users"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

func InitControllers(e *echo.Echo, db *sqlx.DB) {
	echox.LoadControllers(e, []echox.ControllersGetter{
		users.NewModule(db),
	})
}
