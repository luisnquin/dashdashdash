package core

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/core/host/docker"
	"github.com/luisnquin/dashdashdash/internal/core/host/systemd"
	"github.com/luisnquin/dashdashdash/internal/core/users"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

func Init(_ context.Context, e *echo.Echo, db *sqlx.DB) ([]io.Closer, error) {
	usersModule := users.NewModule(db)
	systemdModule := systemd.NewModule()

	dockerModule, err := docker.NewModule()
	if err != nil {
		return nil, err
	}

	echox.LoadControllers(e, []echox.ControllersGetter{
		usersModule, systemdModule, dockerModule,
	})

	return []io.Closer{
		dockerModule,
	}, nil
}
