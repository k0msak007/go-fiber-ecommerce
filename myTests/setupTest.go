package myTests

import (
	"encoding/json"

	"github.com/k0msak007/go-fiber-ecommerce/config"
	"github.com/k0msak007/go-fiber-ecommerce/module/servers"
	"github.com/k0msak007/go-fiber-ecommerce/pkg/databases"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnect(cfg.Db())

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
