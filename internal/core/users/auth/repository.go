package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/luisnquin/dashdashdash/internal/models"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	cache *redis.Client
	db    *sqlx.DB
}

func NewRepository(db *sqlx.DB, cache *redis.Client) Repository { return Repository{cache, db} }

func (r Repository) FindOneUserByUsername(ctx context.Context, username string) (models.User, error) {
	const query = /* sql */ `
	SELECT * FROM users WHERE username = ? LIMIT 1;
	`

	row := r.db.QueryRowxContext(ctx, query, username)

	var user models.User

	if err := row.StructScan(&user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r Repository) SaveUserSession(ctx context.Context, username, token string, expiration time.Duration) error {
	result := r.cache.Set(ctx, fmt.Sprintf("session:%s", username), token, expiration)

	return result.Err()
}
