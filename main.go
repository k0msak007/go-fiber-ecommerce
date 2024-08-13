package main

import (
	"os"

	"github.com/k0msak007/go-fiber-ecommerce/config"
	"github.com/k0msak007/go-fiber-ecommerce/module/servers"
	"github.com/k0msak007/go-fiber-ecommerce/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
