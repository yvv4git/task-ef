package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"
)

type RepoETH struct {
	db *pgxpool.Pool
}

func NewRepoETH(db *pgxpool.Pool) *RepoETH {
	return &RepoETH{
		db: db,
	}
}

type (
	SaveBlockTransactionsETHParams struct {
		Items []TransactionETH
	}

	TransactionETH struct {
		ID               uuid.UUID
		BlockHash        string
		BlockNumber      uint64
		FromAddress      string
		Gas              uint64
		GasPrice         string
		Hash             string
		Input            string
		Nonce            uint64
		ToAddress        string
		TransactionIndex uint
		Value            string
		Type             int
		V                string
		R                string
		S                string
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}
)

func (r *RepoETH) SaveBlockTransactions(ctx context.Context, params SaveBlockTransactionsETHParams) error {
	qb := sq.Insert("eth_transactions").
		Columns("block_hash", "block_number", "from_address", "gas", "gas_price", "hash", "input",
			"nonce", "to_address", "transaction_index", "value", "type", "v", "r", "s")

	for _, item := range params.Items {
		qb = qb.Values(
			item.BlockHash, item.BlockNumber, item.FromAddress, item.Gas, item.GasPrice, item.Hash, item.Input,
			item.Nonce, item.ToAddress, item.TransactionIndex, item.Value, item.Type, item.V, item.R, item.S,
		)
	}

	stmt, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.Exec(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *RepoETH) LastBlockNumber(ctx context.Context) (uint64, error) {
	qb := sq.Select("block_number").
		From("eth_transactions").
		OrderBy("block_number DESC").
		Limit(1)

	stmt, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	var blockNumber uint64
	err = r.db.QueryRow(ctx, stmt, args...).Scan(&blockNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("execute query: %w", err)
	}

	return blockNumber, nil
}

func (r *RepoETH) AddrWithMaxBalacneDiff(ctx context.Context, n uint32) (string, error) {
	const query = `
	WITH
	RecentBlocks AS (
		SELECT
		block_number
		FROM
		eth_transactions
		ORDER BY
		block_number DESC
		LIMIT
		$1
	),
	AddressChanges AS (
		SELECT
		from_address AS address,
		SUM(value) AS change
		FROM
		eth_transactions
		WHERE
		block_number IN (
			SELECT
			block_number
			FROM
			RecentBlocks
		)
		GROUP BY
		from_address
		UNION ALL
		SELECT
		to_address AS address,
		-SUM(value) AS change
		FROM
		eth_transactions
		WHERE
		block_number IN (
			SELECT
			block_number
			FROM
			RecentBlocks
		)
		GROUP BY
		to_address
	),
	AddressBalanceChanges AS (
		SELECT
		address,
		SUM(change) AS balance_change
		FROM
		AddressChanges
		GROUP BY
		address
	)
	SELECT
	address
	FROM
	AddressBalanceChanges
	ORDER BY
	ABS(balance_change) DESC
	LIMIT
	1;
	`

	var result sql.NullString
	if err := r.db.QueryRow(ctx, query, n).Scan(&result); err != nil {
		return "", fmt.Errorf("query: %w", err)
	}

	if !result.Valid {
		return "", ErrInvalidAddrResult
	}

	return result.String, nil
}
