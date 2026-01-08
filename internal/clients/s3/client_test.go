package s3

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	_               = godotenv.Load("../../../.env")
	region          = os.Getenv("AWS_REGION")
	bucket          = os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	accessKeyId     = os.Getenv("AWS_S3_BUCKET_NAME")
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	cfg := Config{
		Region:          region,
		BucketName:      bucket,
		SecretAccessKey: secretAccessKey,
		AccessKeyID:     accessKeyId,
	}
	t.Logf("config: %+v", cfg)
	cli, err := NewClient(ctx, cfg)
	require.NoError(t, err)
	url, err := cli.GeneratePresignedUploadURL(ctx, "bagsy-notion", "application/pdf", time.Minute*10)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	t.Log(url)
}
