package middlewaresRepositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/k0msak007/go-fiber-ecommerce/module/middlewares"
)

type IMiddlewareRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN true ELSE false END)
		FROM "oauth"
		WHERE "user_id" = $1
		AND "access_token" = $2
	`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}

	return true
}

func (r *middlewaresRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
		SELECT
			"id",
			"title"
		FROM "roles"
		ORDER BY "id" DESC;
	`

	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}

	return roles, nil
}
