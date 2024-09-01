package cmd

import (
	"log"
	"red-db-test/database"
	"red-db-test/model"
	"red-db-test/utils"
)

func NewDB() model.DB {
	switch dbName {
	case "RediDB":
		cfg := utils.Config.Database.RediDB
		return utils.Must(database.NewRediDB(cfg.Host, cfg.Port, cfg.Login, cfg.Password, cfg.Database))
	case "MongoDB":
		cfg := utils.Config.Database.MongoDB
		return utils.Must(database.NewMongo(cfg.Endpoint, cfg.AuthSource, cfg.Login, cfg.Password, cfg.Database))
	default:
		log.Fatalf("unknown db: %q", dbName)
		return nil
	}
}
