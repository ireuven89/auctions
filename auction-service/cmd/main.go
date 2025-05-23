package main

import (
	"fmt"

	"github.com/ireuven89/auctions/auction-service/db"
	"github.com/ireuven89/auctions/auction-service/internal"
	"github.com/ireuven89/auctions/shared/config"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.LoadConfig()

	if err != nil {
		panic(fmt.Errorf("failed loading config %v", err))
	}
	logger, _ := zap.NewDevelopment()
	dbConn, err := db.MustNewDB(cfg.Sql.Host, cfg.Sql.User, cfg.Sql.Password, cfg.Sql.Port)

	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	repo := db.NewRepository(dbConn, logger)
	service := internal.NewService(repo, logger)
	transport := internal.NewTransport(service, router)

	transport.ListenAndServe(cfg.Server.Port)
}
