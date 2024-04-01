package config

import (
	"os"
	"strings"
	"time"

	"github.com/luisnquin/go-log"
	"github.com/xlzd/gotp"
)

type Config struct {
	Database Database
	Cache    Cache
	Auth     Auth
}

func New() *Config { return &Config{} }

func (Config) GetIssuerName() string {
	return "dash-dash-dash"
}

type Auth struct{}

func (Auth) GetOPTSecret() string {
	s := mustEnv("OPT_SECRET")
	// gotp.RandomSecret(16)
	if !gotp.IsSecretValid(s) {
		panic("OPT_SECRET is not valid")
	}

	return s
}

func (Auth) GetJWTSecret() []byte {
	return []byte(mustEnv("JWT_SECRET"))
}

func (Auth) GetJWTDuration() time.Duration {
	if v := os.Getenv("JWT_DURATION"); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}

	return time.Hour
}

type Cache struct{}

func (Cache) GetRedisTrustedURL() string {
	return mustEnv("REDIS_TRUSTED_URL")
}

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
