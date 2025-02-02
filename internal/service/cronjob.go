package service

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/domain/usecases"
)

type CronJob struct {
	cfg            *config.Config
	logger         *slog.Logger
	ucETHCollector *usecases.CollectTransactionsETH
	scheduler      *gocron.Scheduler
	running        int32
}

func NewCronJob(cfg *config.Config, logger *slog.Logger, ucETHCollector *usecases.CollectTransactionsETH) *CronJob {
	scheduler := gocron.NewScheduler(time.UTC)

	return &CronJob{
		cfg:            cfg,
		logger:         logger,
		ucETHCollector: ucETHCollector,
		scheduler:      scheduler,
		running:        0,
	}
}

func (c *CronJob) RunAsync(ctx context.Context) {
	_, err := c.scheduler.Every(c.cfg.Cron.Interval).Minutes().Do(func() {
		c.logger.Info("Run job")

		if atomic.CompareAndSwapInt32(&c.running, 0, 1) {
			defer atomic.StoreInt32(&c.running, 0)

			if err := c.ucETHCollector.Process(ctx); err != nil {
				c.logger.Error("Failed to collect transactions", "err", err)
			}
		} else {
			c.logger.Debug("Previous job is still running, skipping new job")
		}
	})
	if err != nil {
		c.logger.Error("Failed to add job to scheduler", "err", err)
	}

	c.scheduler.StartAsync()
}

func (c *CronJob) Stop() {
	c.scheduler.Stop()
}
