package main

import (
	"log"

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
	cfg := config.NewConfig()
	logger := loggerSetup(cfg)
	auther := middleware.NewJwtKey(cfg.TokenSecret)

	locationsStore := locations.NewMemStore(make(map[interface{}]locations.LocationStoreModel))
	locationsSvc := locations.NewDefaultService(logger, locationsStore)

	playersStore := players.NewMemStore(make(map[interface{}]*players.ClientStoreModel))
	playersSvc := players.NewDefaultService(logger, playersStore, cfg.DBTimeOut, cfg.TokenSecret)

	ctlr := app.NewController(logger, locationsSvc, playersSvc, auther)
	worker := pkg.NewWorker(ctlr)
	s.AddWorker("chi-worker", worker)
	s.Run()
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
