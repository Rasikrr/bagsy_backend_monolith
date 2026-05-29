package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetOverview — сводка аналитики (Owner/Manager). Staff → ErrAccessDenied.
func (uc *UseCase) GetOverview(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.OverviewReport, error) {
	if err := requireManagerScope(orgCtx); err != nil {
		return nil, err
	}
	cur, prev, _, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}
	base := analyticsRepo.Scope{
		OrgID:       orgCtx.Organization.ID,
		LocationIDs: resolveLocationIDs(orgCtx, locationID),
	}
	return uc.buildOverview(ctx, base, cur, prev)
}
