package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/luisnquin/go-log"
	"github.com/samber/lo"
)

type Repository struct {
	client *docker.Client
}

func NewRepository(client *docker.Client) Repository { return Repository{client} }

func (r Repository) ListContainers(ctx context.Context, status ...string) ([]models.DockerContainer, error) {
	// https://docs.docker.com/engine/api/v1.29/#tag/Container/operation/ContainerList
	filters := filters.NewArgs()

	if len(status) != 0 && status[0] != "" {
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
		var stdoutStr, stderrStr string

		stdout, err := r.client.ContainerLogs(ctx, cont.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: false,
		})
		if err != nil {
			log.Err(err).Str("container_id", cont.ID).Msg("cannot get container stdout logs")
		} else {
			stdoutStr = mustReadAll(stdout)
			stdout.Close()
		}

		stderr, err := r.client.ContainerLogs(ctx, cont.ID, container.LogsOptions{
			ShowStdout: false,
			ShowStderr: true,
		})
		if err != nil {
			log.Err(err).Str("container_id", cont.ID).Msg("cannot get container stderr logs")
		} else {
			stderrStr = mustReadAll(stderr)
			stderr.Close()
		}

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
			Logs:            models.Stdio{Stdout: stdoutStr, Stderr: stderrStr},
		}
	}), nil
}
