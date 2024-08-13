package databases

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/k0msak007/go-fiber-ecommerce/config"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect db failed: %v\n", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns())

	return db
}
