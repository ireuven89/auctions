package main

import (
	"github.com/ireuven89/auctions/auth-service/db"
	"github.com/ireuven89/auctions/auth-service/internal"
	"github.com/ireuven89/auctions/shared/config"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Color
	loggerCfg.DisableStacktrace = true
	logger, err := loggerCfg.Build()

	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	authDB, err := db.MustNewDB(cfg.Sql.Host, cfg.Sql.User, cfg.Sql.Password, cfg.Sql.Port)
	if err != nil {
		panic(err)
	}
	redisDB, err := db.MustNewRedis(cfg.Redis.Host, cfg.Redis.Password)
	if err != nil {
		panic(err)
	}
	authRepo := db.New(logger, authDB, redisDB)

	router := httprouter.New()
	s, err := internal.NewAuthService(logger, authRepo, "")

	if err != nil {
		panic(err)
	}
	transport := internal.NewTransport(router, s)

	transport.ListenAndServe(cfg.Server.Port)
}
