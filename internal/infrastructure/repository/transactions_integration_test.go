package repository_test

import (
	"context"
	"testing"

	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/infrastructure/logger"
	"github.com/yvv4git/task-ef/internal/infrastructure/repository"
	"github.com/yvv4git/task-ef/internal/utils"
)

func TestRepoETH_SaveBlockTransactions(t *testing.T) {
	// Change directory to project root
	require.NoError(t, utils.ChangeDirToProjectRoot("../../"), "Should change directory without errors")

	// Setup logger & config
	logger := logger.SetupDefaultLogger()
	cfg := config.NewConfig(logger)
	assert.NoError(t, cfg.Load(), "Should load config without errors")

	// Setup postgres pool
	db, err := repository.SetupPostgresPool(cfg.DB)
	require.NoError(t, err, "Should setup postgres pool without errors")

	// Create repo
	repo := repository.NewRepoETH(db)

	err = repo.SaveBlockTransactions(context.Background(), repository.SaveBlockTransactionsETHParams{
		Items: []repository.TransactionETH{
			{
				Hash:             "0x1234567890",
				BlockHash:        "0x1234567890",
				BlockNumber:      1234567,
				FromAddress:      "0x1234567890",
				Gas:              100000,
				GasPrice:         "10",
				Input:            "0x1234567890",
				Nonce:            1,
				ToAddress:        "0x1234567890",
				TransactionIndex: 1,
				Value:            "10",
				Type:             1,
				V:                "1",
				R:                "1",
				S:                "1",
			},
		},
	})
	require.NoError(t, err, "Should save block transactions without errors")
}

func TestRepoETH_LastBlockNumber(t *testing.T) {
	const expectedBlockNumber uint64 = 1234567

	// Change directory to project root
	require.NoError(t, utils.ChangeDirToProjectRoot("../../"), "Should change directory without errors")

	// Setup logger & config
	logger := logger.SetupDefaultLogger()
	cfg := config.NewConfig(logger)
	assert.NoError(t, cfg.Load(), "Should load config without errors")

	// Setup postgres pool
	db, err := repository.SetupPostgresPool(cfg.DB)
	require.NoError(t, err, "Should setup postgres pool without errors")

	// Load fixtures
	require.NoError(t, loadFixtures(db, "eth_transactions"), "Should load fixtures without errors")

	// Create repo
	repo := repository.NewRepoETH(db)

	// Get last block number
	blockNumber, err := repo.LastBlockNumber(context.Background())
	require.NoError(t, err, "Should get last block number without errors")
	assert.Equal(t, expectedBlockNumber, blockNumber, "Should get last block number")
}

func TestRepoETH_LastBlockNumber_EmptyTable(t *testing.T) {
	const expectedBlockNumber uint64 = 0

	// Change directory to project root
	require.NoError(t, utils.ChangeDirToProjectRoot("../../"), "Should change directory without errors")

	// Setup logger & config
	logger := logger.SetupDefaultLogger()
	cfg := config.NewConfig(logger)
	assert.NoError(t, cfg.Load(), "Should load config without errors")

	// Setup postgres pool
	db, err := repository.SetupPostgresPool(cfg.DB)
	require.NoError(t, err, "Should setup postgres pool without errors")

	// Load fixtures
	require.NoError(t, loadFixtures(db, "eth_transactions_empty"), "Should load fixtures without errors")

	// Create repo
	repo := repository.NewRepoETH(db)

	// Get last block number
	blockNumber, err := repo.LastBlockNumber(context.Background())
	require.NoError(t, err, "Should get last block number without errors")
	t.Logf("Last block number: %d", blockNumber)
	assert.Equal(t, expectedBlockNumber, blockNumber, "Should get last block number")
}
