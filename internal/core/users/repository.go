package users

import (
	"github.com/jmoiron/sqlx"
	"github.com/luisnquin/dashdashdash/internal/models"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{db}
}

func (r Repository) FindOneByUsername(username string) (models.User, error) {
	const query = /* sql */ `
	SELECT * FROM users WHERE username = ? LIMIT 1;
	`

	row := r.db.QueryRowx(query, username)

	var user models.User

	if err := row.StructScan(&user); err != nil {
		return models.User{}, err
	}

	user.Password = nil

	return user, nil
}
