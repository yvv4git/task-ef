package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoHealth struct {
	db *pgxpool.Pool
}

func NewRepoHealth(db *pgxpool.Pool) *RepoHealth {
	return &RepoHealth{db: db}
}

func (r *RepoHealth) Check(ctx context.Context) error {
	return r.db.Ping(ctx)
}
