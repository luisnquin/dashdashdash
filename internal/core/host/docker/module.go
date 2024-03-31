package docker

import (
	"io"
	"net/http"

	docker "github.com/docker/docker/client"
	"github.com/luisnquin/dashdashdash/internal/helpers/echox"
)

type (
	Module struct {
		repo    moduleRepository
		closers []io.Closer
	}

	moduleRepository struct {
		docker Repository
	}
)

func NewModule() (Module, error) {
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return Module{}, err
	}

	return Module{
		repo: moduleRepository{
			docker: NewRepository(client),
		},
		closers: []io.Closer{client},
	}, nil
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

func (m Module) Close() error {
	for _, closer := range m.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	return nil
}
