package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetLocation — drill-down по локации. Только Owner с планом Network.
func (uc *UseCase) GetLocation(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID, from, to time.Time) (*domainAnalytics.OverviewReport, error) {
	if !orgCtx.Employee.Role.IsOwner() || orgCtx.Plan.Code != billing.PlanCodeNetwork {
		return nil, domainAnalytics.ErrAccessDenied
	}

	belongs, err := uc.repo.LocationBelongsToOrg(ctx, locationID, orgCtx.Organization.ID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, domainAnalytics.ErrNotFound
	}

	cur, prev, _, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}
	base := analyticsRepo.Scope{
		OrgID:       orgCtx.Organization.ID,
		LocationIDs: []uuid.UUID{locationID},
	}
	return uc.buildOverview(ctx, base, cur, prev)
}
