package users

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

type (
	Module struct {
		repo moduleRepository
	}

	moduleRepository struct {
		users Repository
	}
)

func NewModule(db *sqlx.DB) Module {
	return Module{
		repo: moduleRepository{
			users: NewRepository(db),
		},
	}
}

func (m Module) GetControllers() []echox.Controller {
	return []echox.Controller{
		{
			Method:  http.MethodGet,
			Path:    "/user/:username",
			Handler: m.GetUserByUsernameHandler(),
		},
	}
}
