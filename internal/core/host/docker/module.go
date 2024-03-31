package docker

import (
	"net/http"

	docker "github.com/docker/docker/client"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

type (
	Module struct {
		repo moduleRepository
	}

	moduleRepository struct {
		docker Repository
	}
)

func NewModule() Module {
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	return Module{
		repo: moduleRepository{
			docker: NewRepository(client),
		},
	}
}

func (m Module) GetControllers() []echox.Controller {
	return []echox.Controller{
		{
			Method:  http.MethodGet,
			Path:    "/host/docker/containers",
			Handler: m.GetContainersHandler(),
		},
	}
}
