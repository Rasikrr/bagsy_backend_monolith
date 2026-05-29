package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetClients — аналитика клиентов (Owner/Manager). Staff → ErrAccessDenied.
func (uc *UseCase) GetClients(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.ClientsReport, error) {
	if err := requireManagerScope(orgCtx); err != nil {
		return nil, err
	}
	cur, prev, now, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}

	locs := resolveLocationIDs(orgCtx, locationID)

	statsRows, err := uc.repo.CustomerStats(ctx, orgCtx.Organization.ID, locs, nil)
	if err != nil {
		return nil, err
	}
	stats := toCustomerStats(statsRows)

	base := analyticsRepo.Scope{OrgID: orgCtx.Organization.ID, LocationIDs: locs}
	curVisited, err := uc.repo.CustomersVisited(ctx, scopeForPeriod(base, cur))
	if err != nil {
		return nil, err
	}
	prevVisited, err := uc.repo.CustomersVisited(ctx, scopeForPeriod(base, prev))
	if err != nil {
		return nil, err
	}

	return &domainAnalytics.ClientsReport{
		KPI:       domainAnalytics.NewClientKPISet(stats, toSet(curVisited), toSet(prevVisited), cur, prev, now),
		Segments:  domainAnalytics.NewSegmentBreakdown(stats, now),
		Retention: domainAnalytics.NewRetention(stats),
		Cohorts:   domainAnalytics.NewCohorts(stats, now, cohortMonths),
	}, nil
}
