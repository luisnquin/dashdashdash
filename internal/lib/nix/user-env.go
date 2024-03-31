// Query user environment data. This client is depending on the `nix-env` command.
package nix

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type UserEnvClient struct{}

func NewUserEnvClient() (UserEnvClient, error) {
	p, err := exec.LookPath(NixEnv)
	if err != nil {
		return UserEnvClient{}, err
	} else if p == "" {
		return UserEnvClient{}, fmt.Errorf("couldn't find '%s' command", NixEnv)
	}

	return UserEnvClient{}, nil
}

func (c UserEnvClient) ListInstalledPackages(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, NixEnv, "-q")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSuffix(string(out), "\n"), "\n"), nil
}
