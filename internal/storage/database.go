package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/luisnquin/dashdashdash/internal/config"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func ConnectToTursoDB(config *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("libsql", fmt.Sprintf("%s?authToken=%s", config.Database.TursoDBURL(), config.Database.TursoDBToken()))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}
