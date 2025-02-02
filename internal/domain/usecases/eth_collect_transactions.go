package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yvv4git/task-ef/internal/config"
	"github.com/yvv4git/task-ef/internal/infrastructure/clients"
	"github.com/yvv4git/task-ef/internal/infrastructure/repository"
	"github.com/yvv4git/task-ef/internal/utils"
)

type ExternalClientAPI interface {
	FetchLastBlockNumber() (uint64, error)
	FetchBlockByNumber(blockNumber uint64) (*clients.Block, error)
}

type RepositoryETH interface {
	SaveBlockTransactions(ctx context.Context, params repository.SaveBlockTransactionsETHParams) error
	LastBlockNumber(ctx context.Context) (uint64, error)
}

type CollectTransactionsETH struct {
	cfg       *config.Config
	logger    *slog.Logger
	extClient ExternalClientAPI
	repoETH   RepositoryETH
}

func NewCollectTransactions(cfg *config.Config, log *slog.Logger, extClient ExternalClientAPI, repoETH RepositoryETH) *CollectTransactionsETH {
	return &CollectTransactionsETH{
		cfg:       cfg,
		logger:    log,
		extClient: extClient,
		repoETH:   repoETH,
	}
}

func (ct *CollectTransactionsETH) Process(ctx context.Context) error {
	lastBlockNumber, err := ct.extClient.FetchLastBlockNumber()
	if err != nil {
		return fmt.Errorf("fetch last block number: %w", err)
	}

	processedBlockNumber, err := ct.repoETH.LastBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("get last block number from repository: %w", err)
	}

	if processedBlockNumber == 0 {
		ct.logger.Info("No blocks processed yet")
		if ct.cfg.Cron.StartBlockNumber == 0 {
			return fmt.Errorf("start block number is not set")
		}

		processedBlockNumber = ct.cfg.Cron.StartBlockNumber - 1
	}

	if lastBlockNumber > processedBlockNumber {
		ct.logger.Debug("We must process blocks:", "blocks count", lastBlockNumber-processedBlockNumber)
		for i := processedBlockNumber + 1; i <= lastBlockNumber; i++ {
			block, err := ct.extClient.FetchBlockByNumber(i)
			if err != nil {
				return fmt.Errorf("fetch block by number[%d]: %w", i, err)
			}

			if len(block.Transactions) == 0 {
				ct.logger.Debug("Block has no transactions", "block_number", i)
				continue
			}

			transactions := make([]repository.TransactionETH, len(block.Transactions))
			for j, tx := range block.Transactions {
				txGas, err := utils.HexToUint64(tx.Gas)
				if err != nil {
					ct.logger.Error("Parse tx gas: %v", "err", err)
					continue
				}

				txNonce, err := utils.HexToUint64(tx.Nonce)
				if err != nil {
					ct.logger.Error("Parse tx nonce: %v", "err", err)
					continue
				}

				txType, err := strconv.Atoi(strings.TrimPrefix(tx.Type, "0x"))
				if err != nil {
					ct.logger.Error("Parse tx type", "err", err)
					continue
				}

				transactions[j] = repository.TransactionETH{
					ID:          uuid.New(),
					BlockHash:   block.Hash,
					BlockNumber: i,
					FromAddress: tx.From,
					Gas:         txGas,
					GasPrice:    tx.GasPrice,
					Hash:        tx.Hash,
					Input:       tx.Input,
					Nonce:       txNonce,
					ToAddress:   tx.To,
					Value:       tx.Value,
					Type:        txType,
					V:           tx.V,
					R:           tx.R,
					S:           tx.S,
				}
			}

			err = ct.repoETH.SaveBlockTransactions(ctx, repository.SaveBlockTransactionsETHParams{
				Items: transactions,
			})
			if err != nil {
				return fmt.Errorf("save block transactions: %w", err)
			}

			ct.logger.Debug("Block transactions saved", "block_number", i)
			time.Sleep(time.Millisecond * time.Duration(ct.cfg.Cron.IntervalFetchBlockMS))
		}
	}

	return nil
}
