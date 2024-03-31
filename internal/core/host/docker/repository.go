package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/samber/lo"
)

type Repository struct {
	client *docker.Client
}

func NewRepository(client *docker.Client) Repository { return Repository{client} }

func (r Repository) ListContainers(ctx context.Context, status ...string) ([]models.DockerContainer, error) {
	// https://docs.docker.com/engine/api/v1.29/#tag/Container/operation/ContainerList
	filters := filters.NewArgs()

	if len(status) != 0 {
		filters.Add("status", status[0])
	}

	containers, err := r.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	return lo.Map(containers, func(cont types.Container, _ int) models.DockerContainer {
		stdout, err := r.client.ContainerLogs(ctx, cont.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: false,
		})
		if err != nil {
			log.Panic(err)
		}

		defer stdout.Close()

		stderr, err := r.client.ContainerLogs(ctx, cont.ID, container.LogsOptions{
			ShowStdout: false,
			ShowStderr: true,
		})
		if err != nil {
			log.Panic(err)
		}

		defer stderr.Close()

		return models.DockerContainer{
			ID:              cont.ID,
			Names:           cont.Names,
			Image:           cont.Image,
			ImageID:         cont.ImageID,
			Command:         cont.Command,
			Created:         cont.Created,
			Ports:           cont.Ports,
			SizeRw:          cont.SizeRw,
			SizeRootFs:      cont.SizeRootFs,
			Labels:          cont.Labels,
			State:           cont.State,
			Status:          cont.Status,
			NetworkSettings: cont.NetworkSettings,
			Mounts:          cont.Mounts,
			Logs:            models.Stdio{Stdout: mustReadAll(stdout), Stderr: mustReadAll(stderr)},
		}
	}), nil
}
