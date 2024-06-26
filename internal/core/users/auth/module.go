package auth

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/config"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
	"github.com/redis/go-redis/v9"
	"github.com/xlzd/gotp"
)

type (
	Module struct {
		repo   moduleRepository
		config *config.Config
		totp   *gotp.TOTP
	}

	moduleRepository struct {
		auth Repository
	}
)

func NewModule(config *config.Config, db *sqlx.DB, redis *redis.Client) Module {
	return Module{
		repo: moduleRepository{
			auth: NewRepository(db, redis),
		},
		totp:   gotp.NewDefaultTOTP(config.Auth.GetOPTSecret()),
		config: config,
	}
}

func (m Module) GetControllers() []echox.Controller {
	return []echox.Controller{
		{
			Method:  http.MethodGet,
			Path:    "/auth/totp/generate",
			Handler: m.GenerateTOTPUriHandler(), // Auth Basic || JWT -> uri
			Middlewares: []echo.MiddlewareFunc{
				m.AuthCheckMiddleware(),
			},
			Verbose: true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/totp/validate/:code",
			Handler: m.ValidateTOTPCodeHandler(), // Auth Basic || JWT -> TOTP code
			Middlewares: []echo.MiddlewareFunc{
				m.AuthCheckMiddleware(),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: m.LoginHandler(), // Auth Basic -> JWT
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: m.LogoutHandler(), // Auth Basic || JWT -> no JWT
		},
	}
}
