package nix

import (
	"context"

	nix_lib "github.com/luisnquin/dashdashdash/internal/lib/nix"
)

type Repository struct {
	nixEnvClient nix_lib.UserEnvClient
}

func NewRepository() (Repository, error) {
	nixEnvClient, err := nix_lib.NewUserEnvClient()
	if err != nil {
		return Repository{}, err
	}

	return Repository{nixEnvClient}, nil
}

func (r Repository) GetInstalledUserEnvPackages(ctx context.Context) ([]string, error) {
	return r.nixEnvClient.ListInstalledPackages(ctx)
}
