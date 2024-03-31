package docker

import (
	"io"
)

func getContainerStatusList() []string {
	return []string{"created", "restarting", "running", "removing", "paused", "exited", "dead"}
}

func mustReadAll(r io.Reader) string {
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return string(b)
}
