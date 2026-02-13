package billing

import (
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Plan
// ─────────────────────────────────────────────────────────────────

type Plan struct {
	ID           uuid.UUID
	Slug         shared.Slug // 'solo', 'point', 'network'
	Name         string
	Description  *string
	PriceMonthly shared.Money
	PriceAnnual  shared.Money
	SortOrder    int
	Active       bool
	Capabilities []PlanCapability
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type PlanCapability struct {
	ID         uuid.UUID
	PlanID     uuid.UUID
	Resource   Resource // 'analytics', 'max_locations'
	LimitValue *int     // NULL = unlimited
}

func NewPlan(
	name string,
	description *string,
	priceMonthly, priceAnnual shared.Money,
) (*Plan, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrPlanNameRequired
	}

	slug, err := shared.NewSlug(name)
	if err != nil {
		return nil, err
	}

	return &Plan{
		ID:           uuid.New(),
		Slug:         slug,
		Name:         strings.TrimSpace(name),
		Description:  description,
		PriceMonthly: priceMonthly,
		PriceAnnual:  priceAnnual,
		Active:       true,
		Capabilities: make([]PlanCapability, 0),
		CreatedAt:    time.Now(),
	}, nil
}

func (p *Plan) AddCapability(resource Resource, limit *int) {
	capability := PlanCapability{
		ID:         uuid.New(),
		PlanID:     p.ID,
		Resource:   resource,
		LimitValue: limit,
	}
	p.Capabilities = append(p.Capabilities, capability)
	p.touch()
}

func (p *Plan) Deactivate() {
	p.Active = false
	p.touch()
}

func (p *Plan) Activate() {
	p.Active = true
	p.touch()
}

func (p *Plan) touch() {
	now := time.Now()
	p.UpdatedAt = &now
}
