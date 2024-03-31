package config

import (
	"os"
	"strings"

	"github.com/luisnquin/go-log"
)

type Config struct {
	Database
}

func New() *Config { return &Config{} }

type Database struct{}

func (Database) TursoDBURL() string {
	return mustEnv("TURSO_DB_URL")
}

func (Database) TursoDBToken() string {
	return mustEnv("TURSO_DB_TOKEN")
}

func mustEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		log.Error().Msgf("environment variable '%s' is required", key)
	}

	return value
}
