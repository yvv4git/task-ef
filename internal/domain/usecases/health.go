package usecases

import "context"

type HealthRepository interface {
	Check(context.Context) error
}

type Health struct {
	healthRepository HealthRepository
}

func NewHealth(healthRepository HealthRepository) *Health {
	return &Health{
		healthRepository: healthRepository,
	}
}

func (h *Health) Check(ctx context.Context) error {
	return h.healthRepository.Check(ctx)
}
