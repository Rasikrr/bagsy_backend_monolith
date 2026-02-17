package plan

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type planModel struct {
	ID           uuid.UUID       `db:"id"`
	Code         string          `db:"code"`
	Name         string          `db:"name"`
	Description  *string         `db:"description"`
	PriceMonthly decimal.Decimal `db:"price_monthly"`
	PriceAnnual  decimal.Decimal `db:"price_annual"`
	SortOrder    int             `db:"sort_order"`
	Active       bool            `db:"active"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    *time.Time      `db:"updated_at"`
}

type capabilityModel struct {
	ID         uuid.UUID `db:"id"`
	PlanID     uuid.UUID `db:"plan_id"`
	Resource   string    `db:"resource"`
	LimitValue *int      `db:"limit_value"`
}

func (m *planModel) toDomain(caps []capabilityModel) (*billing.Plan, error) {
	code, err := billing.ParsePlanCode(m.Code)
	if err != nil {
		return nil, err
	}

	priceMonthly, err := shared.NewMoney(m.PriceMonthly)
	if err != nil {
		return nil, err
	}

	priceAnnual, err := shared.NewMoney(m.PriceAnnual)
	if err != nil {
		return nil, err
	}

	capabilities := make([]billing.PlanCapability, 0, len(caps))
	for _, c := range caps {
		var limit billing.Limit
		if c.LimitValue == nil {
			limit = billing.NewUnlimited()
		} else {
			limit, err = billing.NewLimit(*c.LimitValue)
			if err != nil {
				return nil, err
			}
		}

		capabilities = append(capabilities, billing.PlanCapability{
			ID:       c.ID,
			PlanID:   c.PlanID,
			Resource: billing.Resource(c.Resource),
			Limit:    limit,
		})
	}

	return &billing.Plan{
		ID:           m.ID,
		Code:         code,
		Name:         m.Name,
		Description:  m.Description,
		PriceMonthly: priceMonthly,
		PriceAnnual:  priceAnnual,
		SortOrder:    m.SortOrder,
		Active:       m.Active,
		Capabilities: capabilities,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}
