package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) UpdateProfile(ctx context.Context, employeeID uuid.UUID, input UpdateProfileInput) (*ProfileOutput, error) {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, errors.Wrap(err, "get employee")
	}

	if err = emp.UpdateProfile(input.FirstName, input.LastName); err != nil {
		return nil, err
	}

	if input.AvatarID != nil {
		asset, assetErr := u.mediaRepo.GetByID(ctx, *input.AvatarID)
		if assetErr != nil {
			return nil, errors.Wrap(assetErr, "get avatar asset")
		}
		if !asset.IsReady() {
			return nil, media.ErrAssetNotReady
		}
		if err = emp.ChangeAvatar(*input.AvatarID); err != nil {
			return nil, err
		}
	}

	if err = u.employeeRepo.Save(ctx, emp); err != nil {
		return nil, errors.Wrap(err, "save employee")
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
