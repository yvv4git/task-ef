package clients_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/infrastructure/clients"
	"github.com/yvv4git/task-ef/internal/infrastructure/logger"
	"github.com/yvv4git/task-ef/internal/utils"
)

func TestClient_FetchLastBlockNumber(t *testing.T) {
	// Change directory to project root
	require.NoError(t, utils.ChangeDirToProjectRoot("../../"), "Should change directory without errors")

	// Setup logger & config
	logger := logger.SetupDefaultLogger()
	cfg := config.NewConfig(logger)
	assert.NoError(t, cfg.Load(), "Should load config without errors")

	// Send request to external API
	extClient := clients.NewClient(cfg)
	lastBlockNum, err := extClient.FetchLastBlockNumber()
	require.NoError(t, err, "Should fetch last block number without errors")
	t.Logf("Last block number: %d", lastBlockNum)
	require.Greaterf(t, lastBlockNum, uint64(0), "Should get positive block number, got: %d", lastBlockNum)
}

func TestClient_FetchLastBlockByNumber(t *testing.T) {
	const blockNumber uint64 = 19000000
	// Change directory to project root
	require.NoError(t, utils.ChangeDirToProjectRoot("../../"), "Should change directory without errors")

	// Setup logger & config
	logger := logger.SetupDefaultLogger()
	cfg := config.NewConfig(logger)
	assert.NoError(t, cfg.Load(), "Should load config without errors")

	// Send request to external API
	extClient := clients.NewClient(cfg)
	block, err := extClient.FetchBlockByNumber(blockNumber)
	require.NoError(t, err, "Should fetch block by number without errors")
	require.NotEmpty(t, block, "Should get block with transactions, got: %v", block)
}
