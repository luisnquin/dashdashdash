package models

import "github.com/docker/docker/api/types"

type DockerContainer struct {
	ID              string                        `json:"id"`
	Names           []string                      `json:"names"`
	Image           string                        `json:"image"`
	ImageID         string                        `json:"imageId"`
	Command         string                        `json:"command"`
	Created         int64                         `json:"created"`
	Ports           []types.Port                  `json:"ports"`
	SizeRw          int64                         `json:"sizeRw"`
	SizeRootFs      int64                         `json:"sizeRootFs"`
	Labels          map[string]string             `json:"labels"`
	State           string                        `json:"state"`
	Status          string                        `json:"status"`
	NetworkSettings *types.SummaryNetworkSettings `json:"networkSettings"`
	Mounts          []types.MountPoint            `json:"mounts"`
	Logs            Stdio                         `json:"logs"`
}
