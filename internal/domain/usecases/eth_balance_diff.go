package usecases

import "context"

type EthBalanceDiffRepository interface {
	AddrWithMaxBalacneDiff(ctx context.Context, n uint32) (string, error)
}

type EthBalanceDiff struct {
	repo EthBalanceDiffRepository
}

func NewEthBalanceDiff(repo EthBalanceDiffRepository) *EthBalanceDiff {
	return &EthBalanceDiff{
		repo: repo,
	}
}

func (e EthBalanceDiff) FindAddrWithMaxBalanceDiff(ctx context.Context, n uint32) (string, error) {
	return e.repo.AddrWithMaxBalacneDiff(ctx, n)
}
