package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/domain/usecases"
	"github.com/yvv4git/task-ef/internal/infrastructure/api"
	"github.com/yvv4git/task-ef/internal/infrastructure/logger"
	"github.com/yvv4git/task-ef/internal/infrastructure/repository"
	"github.com/yvv4git/task-ef/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defaultLogger := logger.SetupDefaultLogger()

	cfg := config.NewConfig(defaultLogger)
	if err := cfg.Load(); err != nil {
		defaultLogger.Error("load config", "err", err)
	}

	log := logger.SetupLoggerWithLevel(cfg.LogLevel)
	log.Info("Service api started")
	defer log.Info("Service api shutdown")

	db, err := repository.SetupPostgresPool(cfg.DB)
	if err != nil {
		log.Error("Setup postgres pool", "err", err)
	}

	repoHealth := repository.NewRepoHealth(db)
	repoETH := repository.NewRepoETH(db)

	ucHealth := usecases.NewHealth(repoHealth)
	ucAnalyzeETH := usecases.NewEthBalanceDiff(repoETH)

	svcHealth := service.NewHealth(ucHealth)
	svcAnalyzeETH := service.NewAnalyzeETH(ucAnalyzeETH)

	webSrv := api.NewWeb(&cfg.API, log, svcHealth, svcAnalyzeETH)
	if err := webSrv.Run(ctx); err != nil {
		log.Error("Run web server", "err", err)
	}
}
