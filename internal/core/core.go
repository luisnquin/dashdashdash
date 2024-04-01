package core

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/config"
	"github.com/luisnquin/dashdashdash/internal/core/host/docker"
	"github.com/luisnquin/dashdashdash/internal/core/host/nix"
	"github.com/luisnquin/dashdashdash/internal/core/host/systemd"
	"github.com/luisnquin/dashdashdash/internal/core/users"
	"github.com/luisnquin/dashdashdash/internal/core/users/auth"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
	"github.com/redis/go-redis/v9"
)

func Init(_ context.Context, e *echo.Echo, config *config.Config, db *sqlx.DB, cache *redis.Client) ([]io.Closer, error) {
	usersModule := users.NewModule(db)
	systemdModule := systemd.NewModule()
	authModule := auth.NewModule(config, db, cache)

	dockerModule, err := docker.NewModule()
	if err != nil {
		return nil, err
	}

	nixModule, err := nix.NewModule()
	if err != nil {
		return nil, err
	}

	echox.LoadControllers(e, []echox.ControllersGetter{
		usersModule, systemdModule, dockerModule, nixModule,
		authModule,
	})

	return []io.Closer{
		dockerModule,
	}, nil
}
