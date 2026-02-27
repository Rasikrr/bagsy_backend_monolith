package media

import (
	"context"

	mediaDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/s3"
	"github.com/cockroachdb/errors"
)

func (u *UseCase) GenerateUploadURL(
	ctx context.Context,
	input GenerateUploadURLInput,
) (*GenerateUploadURLOutput, error) {
	mimeType, err := mediaDomain.ParseMimeType(input.MimeType)
	if err != nil {
		return nil, err
	}

	if input.SizeBytes > u.maxFileSizeBytes {
		return nil, mediaDomain.ErrInvalidFileSize
	}

	purpose, err := mediaDomain.ParsePurpose(input.Purpose)
	if err != nil {
		return nil, err
	}

	asset, err := mediaDomain.NewAsset(mediaDomain.CreateAssetParams{
		Bucket:    u.storage.BucketName(),
		Filename:  input.Filename,
		Purpose:   purpose,
		MimeType:  mimeType,
		SizeBytes: input.SizeBytes,
	})
	if err != nil {
		return nil, err
	}

	if err = u.mediaRepo.Save(ctx, asset); err != nil {
		return nil, errors.Wrap(err, "save media asset")
	}

	policy, err := u.storage.GeneratePresignedPostURL(ctx, s3.UploadPolicyOptions{
		Key:              asset.ObjectKey,
		ContentType:      mimeType.String(),
		ContentLengthMin: 1,
		ContentLengthMax: u.maxFileSizeBytes,
		Expires:          u.uploadExpires,
	})
	if err != nil {
		return nil, errors.Wrap(err, "generate presigned post url")
	}

	return &GenerateUploadURLOutput{
		AssetID:      asset.ID,
		UploadURL:    policy.URL,
		UploadFields: policy.Fields,
	}, nil
}
