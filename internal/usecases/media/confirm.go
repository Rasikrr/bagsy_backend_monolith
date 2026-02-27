package media

import (
	"context"

	mediaDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) ConfirmUpload(ctx context.Context, assetID uuid.UUID) error {
	asset, err := u.mediaRepo.GetByID(ctx, assetID)
	if err != nil {
		return err
	}

	exists, err := u.storage.Exists(ctx, asset.ObjectKey)
	if err != nil {
		return errors.Wrap(err, "check s3 object existence")
	}
	if !exists {
		return mediaDomain.ErrS3ObjectNotFound
	}

	if err = asset.MarkAsUploaded(); err != nil {
		return err
	}

	if err = u.mediaRepo.Save(ctx, asset); err != nil {
		return errors.Wrap(err, "save media asset")
	}

	return nil
}
