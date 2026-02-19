package access

import (
	"encoding/json"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type orgContextModel struct {
	EmployeeID                 uuid.UUID  `db:"employee_id"`
	EmployeePhone              string     `db:"employee_phone"`
	EmployeeLocationID         *uuid.UUID `db:"employee_location_id"`
	EmployeeRole               string     `db:"employee_role"`
	EmployeeCanProvideServices bool       `db:"employee_can_provide_services"`
	EmployeeCanManageSchedule  bool       `db:"employee_can_manage_location_schedule"`

	OrganizationID     uuid.UUID `db:"organization_id"`
	OrganizationActive bool      `db:"organization_active"`

	SubscriptionStatus *string `db:"subscription_status"`

	PlanCode         *string `db:"plan_code"`
	PlanCapabilities []byte  `db:"plan_capabilities"`
}

func (m *orgContextModel) toDomain() (*access.OrgContext, error) {
	phone, err := shared.NewPhone(m.EmployeePhone)
	if err != nil {
		return nil, fmt.Errorf("parse employee phone: %w", err)
	}

	role, err := identity.ParseRole(m.EmployeeRole)
	if err != nil {
		return nil, err
	}

	var locID uuid.UUID
	if m.EmployeeLocationID != nil {
		locID = *m.EmployeeLocationID
	}

	empInfo := access.EmployeeInfo{
		ID:         m.EmployeeID,
		Phone:      phone,
		LocationID: locID,
		Role:       role,
		Permissions: identity.NewPermissions(
			m.EmployeeCanProvideServices,
			m.EmployeeCanManageSchedule,
		),
	}

	orgInfo := access.OrganizationInfo{
		ID:     m.OrganizationID,
		Active: m.OrganizationActive,
	}

	subInfo := access.SubscriptionInfo{}

	if m.SubscriptionStatus != nil {
		subStatus, err := billing.ParseSubscriptionStatus(*m.SubscriptionStatus)
		if err != nil {
			return nil, err
		}
		subInfo.Status = subStatus

	}

	var planCode billing.PlanCode
	if m.PlanCode != nil {
		planCode, err = billing.ParsePlanCode(*m.PlanCode)
		if err != nil {
			return nil, err
		}
	}

	capabilitiesMap := make(map[billing.Resource]billing.Limit)
	if len(m.PlanCapabilities) > 0 && string(m.PlanCapabilities) != "null" {
		var rawCaps map[string]*int
		if err := json.Unmarshal(m.PlanCapabilities, &rawCaps); err != nil {
			return nil, fmt.Errorf("unmarshal plan capabilities: %w", err)
		}

		for resName, val := range rawCaps {
			var limit billing.Limit
			if val == nil {
				limit = billing.NewUnlimited()
			} else {
				// We assume valid values in DB per domain rules
				limit, _ = billing.NewLimit(*val)
			}
			capabilitiesMap[billing.Resource(resName)] = limit
		}
	}

	planInfo := access.PlanInfo{
		Code:         planCode,
		Capabilities: access.NewCapabilities(capabilitiesMap),
	}

	return &access.OrgContext{
		Employee:     empInfo,
		Organization: orgInfo,
		Subscription: subInfo,
		Plan:         planInfo,
	}, nil
}
