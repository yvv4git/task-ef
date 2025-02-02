package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/domain/usecases"
	"github.com/yvv4git/task-ef/internal/infrastructure/clients"
	"github.com/yvv4git/task-ef/internal/infrastructure/logger"
	"github.com/yvv4git/task-ef/internal/infrastructure/repository"
	"github.com/yvv4git/task-ef/internal/service"
	"github.com/yvv4git/task-ef/internal/utils"
)

func main() {
	utils.ChangeDirToProjectRoot("../../")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defaultLogger := logger.SetupDefaultLogger()

	cfg := config.NewConfig(defaultLogger)
	if err := cfg.Load(); err != nil {
		defaultLogger.Error("load config", "err", err)
	}

	log := logger.SetupLoggerWithLevel(cfg.LogLevel)
	log.Info("Service cron-job started")
	defer log.Info("Service cron-job shutdown")

	db, err := repository.SetupPostgresPool(cfg.DB)
	if err != nil {
		log.Error("setup postgres db pool", "err", err)
	}

	repo := repository.NewRepoETH(db)
	extClient := clients.NewClient(cfg)
	uc := usecases.NewCollectTransactions(cfg, log, extClient, repo)
	svc := service.NewCronJob(cfg, log, uc)
	svc.RunAsync(ctx)
	defer svc.Stop()

	<-ctx.Done()
}
