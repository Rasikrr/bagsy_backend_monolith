package billing

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
)

func (u *UseCase) ListPlans(ctx context.Context) ([]*billing.Plan, error) {
	plans, err := u.planRepo.FindAllActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	return plans, nil
}
