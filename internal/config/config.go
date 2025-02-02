package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		logger      *slog.Logger
		DB          DB          `envconfig:"POSTGRES"`
		API         API         `envconfig:"API"`
		ExternalAPI ExternalAPI `envconfig:"EXT_API"`
		Cron        CronJob     `envconfig:"CRON"`
		LogLevel    slog.Level  `envconfig:"LOG_LEVEL"`
	}

	DB struct {
		Host     string `envconfig:"HOST"`
		Port     int    `envconfig:"PORT"`
		Name     string `envconfig:"DB"`
		User     string `envconfig:"USER"`
		Password string `envconfig:"PASSWORD"`
		SSLMode  string `envconfig:"SSL_MODE"`
	}

	API struct {
		Host              string `envconfig:"HOST"`
		Port              int    `envconfig:"PORT"`
		ShutdownTimeout   int    `envconfig:"SHUTDOWN_TIMEOUT_SECONDS"`
		AnalyzeTxCountETH uint32 `envconfig:"ANALYZE_TX_COUNT_ETH"`
	}

	ExternalAPI struct {
		Key     string `envconfig:"KEY"`
		URL     string `envconfig:"URL"`
		Timeout int    `envconfig:"TIMEOUT_SECONDS"`
	}

	CronJob struct {
		Interval             int    `envconfig:"INTERVAL_MINUTES"`
		StartBlockNumber     uint64 `envconfig:"START_BLOCK_NUMBER"`
		IntervalFetchBlockMS int    `envconfig:"INTERVAL_FETCH_BLOCK_MS"`
	}
)

func NewConfig(logger *slog.Logger) *Config {
	return &Config{
		logger: logger,
	}
}

func (c *Config) Load() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		c.logger.Debug("failed to load .env file", "error", err)
	}

	// Parse environment variables
	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}
