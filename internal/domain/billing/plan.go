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
	Code         PlanCode
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
	ID       uuid.UUID
	PlanID   uuid.UUID
	Resource Resource // 'analytics', 'max_locations'
	Limit    Limit
}

func NewPlan(
	name string,
	description *string,
	planCode PlanCode,
	priceMonthly, priceAnnual shared.Money,
) (*Plan, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrPlanNameRequired
	}

	return &Plan{
		ID:           uuid.New(),
		Code:         planCode,
		Name:         strings.TrimSpace(name),
		Description:  description,
		PriceMonthly: priceMonthly,
		PriceAnnual:  priceAnnual,
		Active:       true,
		Capabilities: make([]PlanCapability, 0),
		CreatedAt:    time.Now(),
	}, nil
}

func (p *Plan) AddCapability(resource Resource, limit Limit) {
	for i, capability := range p.Capabilities {
		if capability.Resource == resource {
			// If already exists, update the limit
			p.Capabilities[i].Limit = limit
			p.touch()
			return
		}
	}

	// If doesn't exist, append new one
	capability := PlanCapability{
		ID:       uuid.New(),
		PlanID:   p.ID,
		Resource: resource,
		Limit:    limit,
	}
	p.Capabilities = append(p.Capabilities, capability)
	p.touch()
}

func (p *Plan) ChangeCapabilityLimit(resource Resource, limit Limit) error {
	for i, capability := range p.Capabilities {
		if capability.Resource == resource {
			p.Capabilities[i].Limit = limit
			p.touch()
			return nil
		}
	}
	return ErrPlanCapabilityNotFound
}

func (p *Plan) RemoveCapability(resource Resource) {
	for i, capability := range p.Capabilities {
		if capability.Resource == resource {
			p.Capabilities = append(p.Capabilities[:i], p.Capabilities[i+1:]...)
			p.touch()
			return
		}
	}
	return
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
