package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/infrastructure/logger"
	"github.com/yvv4git/task-ef/internal/utils"
)

func TestConfigLoad(t *testing.T) {
	// Change directory to project root
	require.NoError(t, utils.ChangeDirToProjectRoot("../../"), "Should change directory without errors")

	// Setup logger & config
	logger := logger.SetupDefaultLogger()
	cfg := config.NewConfig(logger)

	// Load config
	assert.NoError(t, cfg.Load(), "Should load config without errors")
}
