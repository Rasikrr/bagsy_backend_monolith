package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetFinance — финансовый отчёт (Owner/Manager). Staff → ErrAccessDenied.
func (uc *UseCase) GetFinance(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.FinanceReport, error) {
	if err := requireManagerScope(orgCtx); err != nil {
		return nil, err
	}
	cur, _, _, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}

	scope := analyticsRepo.Scope{
		OrgID:       orgCtx.Organization.ID,
		LocationIDs: resolveLocationIDs(orgCtx, locationID),
		From:        cur.From,
		To:          cur.To,
	}

	finRows, err := uc.repo.EmployeeFinance(ctx, scope)
	if err != nil {
		return nil, err
	}

	inputs := make([]domainAnalytics.PayrollInput, 0, len(finRows))
	for _, r := range finRows {
		inputs = append(inputs, domainAnalytics.PayrollInput{
			EmployeeID:        r.EmployeeID,
			FullName:          r.FullName,
			CommissionPercent: r.CommissionPercent,
			Revenue:           r.Revenue,
		})
	}

	report := domainAnalytics.NewFinanceReport(inputs)
	return &report, nil
}
