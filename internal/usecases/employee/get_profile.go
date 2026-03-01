package employee

import (
	"context"

	"github.com/google/uuid"
)

func (u *UseCase) GetProfile(ctx context.Context, employeeID uuid.UUID) (*ProfileOutput, error) {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	avatarURL, err := u.resolveAvatarURL(ctx, emp.AvatarID)
	if err != nil {
		return nil, err
	}

	return &ProfileOutput{
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
	}, nil
}
