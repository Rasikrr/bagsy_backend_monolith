package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) GetProfile(ctx context.Context, employeeID uuid.UUID) (*ProfileOutput, error) {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	var avatarURL *string
	if emp.AvatarID != nil {
		var asset *media.Asset
		asset, err = u.mediaRepo.GetByID(ctx, *emp.AvatarID)
		if err != nil {
			return nil, errors.Wrap(err, "get avatar asset")
		}

		if asset.IsReady() {
			var url string
			url, err = u.storage.GeneratePresignedDownloadURL(ctx, asset.ObjectKey, u.avatarURLExpiry)
			if err != nil {
				return nil, errors.Wrap(err, "generate avatar url")
			}
			avatarURL = &url
		}
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
