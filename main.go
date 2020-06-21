package main

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/mattn/go-colorable"
	"github.com/voi-oss/svc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"geogame/config"
	"geogame/internal/app"
	"geogame/internal/locations"
	"geogame/internal/middleware"
	"geogame/internal/players"
	"geogame/pkg"
)

func main() {
	s, err := svc.New("geogame", "snap-shot")
	svc.MustInit(s, err)

	// config setup
	cfg := config.NewConfig()

	// logger setup
	logger := loggerSetup(cfg)

	// middleware authentication setup
	auther := middleware.NewJwtKey(cfg.TokenSecret)

	// postgres setup
	pgConfig := &pkg.Config{}
	svc.MustInit(s, svc.LoadFromEnv(pgConfig))

	pgWorker := pkg.New("pg", pkg.WithConfig(pgConfig))
	svc.MustInit(s, pgWorker.Connect())

	// setup locations service
	locationsStore := newLocationsStore(cfg, pgWorker.DB(), logger)
	locationsSvc := locations.NewDefaultService(logger, locationsStore)

	// setup players service
	playersStore := newPlayersStore(cfg, pgWorker.DB(), logger)
	playersSvc := players.NewDefaultService(logger, playersStore, cfg.DBTimeOut, cfg.TokenSecret)

	// init controller
	controller := app.NewController(logger, locationsSvc, playersSvc, auther)
	HTTPWorker := pkg.NewChiWorker(controller)

	s.AddWorker("pg-worker", pgWorker)
	s.AddWorker("http-worker", HTTPWorker)
	s.Run()
}

func newPlayersStore(cfg *config.Config, db *sqlx.DB, logger *zap.Logger) players.Store {
	if cfg.Env == config.EnvDev {
		return players.NewMemStore(make(map[interface{}]*players.ClientStoreModel))
	}
	return players.NewPostgres(db, logger)
}

func newLocationsStore(cfg *config.Config, db *sqlx.DB, logger *zap.Logger) locations.Store {
	if cfg.Env == config.EnvDev {
		return locations.NewMemStore(make(map[interface{}]locations.LocationStoreModel))
	}
	return locations.NewPostgres(db, logger)
}

func loggerSetup(cfg *config.Config) *zap.Logger {
	if cfg.Env == config.EnvProd {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("failed to create zap logger : %v", err)
		}
		return logger
	}
	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder
	bb := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	bb.Warn("logger setup done")
	return bb
}
