package service

import "context"

type HealthUseCase interface {
	Check(ctx context.Context) error
}

type Health struct {
	ucHealth HealthUseCase
}

func NewHealth(ucHealth HealthUseCase) *Health {
	return &Health{
		ucHealth: ucHealth,
	}
}

func (h *Health) Check(ctx context.Context) error {
	return h.ucHealth.Check(ctx)
}
