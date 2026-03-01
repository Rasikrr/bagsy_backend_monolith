package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) GetList(ctx context.Context, orgCtx *access.OrgContext, filter *identity.EmployeeFilter) (*ListOutput, error) {
	filter.OrganizationID = orgCtx.Organization.ID

	if err := u.policy.CanListEmployees(orgCtx, filter); err != nil {
		return nil, err
	}

	page, err := u.employeeRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "get employees by filter")
	}

	avatarIDs := make([]uuid.UUID, 0)
	for _, emp := range page.Items {
		if emp.AvatarID != nil {
			avatarIDs = append(avatarIDs, *emp.AvatarID)
		}
	}

	avatarURLs, err := u.resolveAvatarURLsBatch(ctx, avatarIDs)
	if err != nil {
		return nil, err
	}

	items := make([]ProfileOutput, 0, len(page.Items))
	for _, emp := range page.Items {
		var avatarURL *string
		if emp.AvatarID != nil {
			if url, ok := avatarURLs[*emp.AvatarID]; ok {
				avatarURL = &url
			}
		}

		items = append(items, ProfileOutput{
			ID:             emp.ID,
			Phone:          emp.Phone.String(),
			FirstName:      emp.FirstName,
			LastName:       emp.LastName,
			AvatarURL:      avatarURL,
			OrganizationID: emp.OrganizationID,
			LocationID:     emp.LocationID,
			Role:           emp.Role,
			Permissions:    emp.Permissions,
			Active:         emp.Active,
			CreatedAt:      emp.CreatedAt,
		})
	}

	return &ListOutput{
		Items: items,
		Total: page.Total,
	}, nil
}
