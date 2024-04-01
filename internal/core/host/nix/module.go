package nix

import (
	"net/http"

	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

type (
	Module struct {
		repo moduleRepository
	}

	moduleRepository struct {
		nix Repository
	}
)

func NewModule() (Module, error) {
	nixRepo, err := NewRepository()
	if err != nil {
		return Module{}, err
	}

	return Module{
		repo: moduleRepository{
			nix: nixRepo,
		},
	}, nil
}

func (m Module) GetControllers() []echox.Controller {
	return []echox.Controller{
		{
			Method:  http.MethodGet,
			Path:    "/host/nix/user-env/packages",
			Handler: m.GetInstalledUserEnvPackagesHandler(),
			Auth:    true,
		},
	}
}
