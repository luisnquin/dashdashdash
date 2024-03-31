package docker

import (
	"io"
	"log"
)

func getContainerStatusList() []string {
	return []string{"created", "restarting", "running", "removing", "paused", "exited", "dead"}
}

func mustReadAll(r io.Reader) string {
	b, err := io.ReadAll(r)
	if err != nil {
		log.Panic(err)
	}

	return string(b)
}
