package service

import (
	"context"

	"github.com/yvv4git/task-ef/internal/domain/usecases"
)

type AnalyzeETH struct {
	ucETHBalanceDiff *usecases.EthBalanceDiff
}

func NewAnalyzeETH(ucETHBalanceDiff *usecases.EthBalanceDiff) *AnalyzeETH {
	return &AnalyzeETH{
		ucETHBalanceDiff: ucETHBalanceDiff,
	}
}

func (e *AnalyzeETH) FindAddressWithMaxDiffBalance(ctx context.Context, n uint32) (string, error) {
	return e.ucETHBalanceDiff.FindAddrWithMaxBalanceDiff(ctx, n)
}
