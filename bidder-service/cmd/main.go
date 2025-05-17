package main

import (
	"fmt"
	"github.com/ireuven89/auctions/bidder-service/db"
	"github.com/ireuven89/auctions/bidder-service/internal"
	"github.com/ireuven89/auctions/shared/config"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(fmt.Errorf("failed loading config %v", err))
	}
	dbConn, err := db.MustNewDB(cfg.Sql.Host, cfg.Sql.User, cfg.Sql.Password, cfg.Sql.Port)

	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()
	repo := db.NewRepository(dbConn, logger)
	router := httprouter.New()

	service := internal.NewService(repo, logger)
	transport := internal.NewTransport(router, service)

	transport.ListenAndServe("8090")
}
