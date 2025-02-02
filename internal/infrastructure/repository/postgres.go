package repository

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yvv4git/task-ef/internal/config"
)

const (
	SSLMode = "require"
)

func SetupPostgresPool(dbCfg config.DB) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
		dbCfg.SSLMode,
	)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.MinConns = 1
	cfg.MaxConns = 20
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	cfg.ConnConfig.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
		"timezone":                    "UTC",
	}

	if dbCfg.SSLMode == SSLMode {
		cfg.ConnConfig.TLSConfig = &tls.Config{} // TODO: if you want use ssl
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return dbPool, nil
}
