package s3

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	_               = godotenv.Load("../../.env")
	region          = os.Getenv("AWS_REGION")
	endpoint        = os.Getenv("AWS_S3_ENDPOINT")
	bucket          = os.Getenv("AWS_S3_BUCKET_NAME")
	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	accessKeyId     = os.Getenv("AWS_ACCESS_KEY_ID")
)

type S3TestSuite struct {
	suite.Suite
	cli *Client
}

// SetupSuite запускается ОДИН РАЗ перед выполнением всех тестов в S3TestSuite
func (s *S3TestSuite) SetupSuite() {
	s.T().Skip()
	ctx := context.Background()

	cfg := Config{
		Region:          region,
		BucketName:      bucket,
		Endpoint:        endpoint,
		SecretAccessKey: secretAccessKey,
		AccessKeyID:     accessKeyId,
	}

	cli, err := NewClient(ctx, cfg)
	s.Require().NoError(err, "Не удалось инициализировать S3 клиента")
	s.cli = cli
}

// Теперь тесты — это методы нашей структуры

func (s *S3TestSuite) TestGenPresignedUploadURL() {
	ctx := context.Background()
	url, err := s.cli.GeneratePresignedUploadURL(ctx, "bagsy-notion", "application/pdf", time.Minute*10)

	s.Require().NoError(err)
	s.Require().NotEmpty(url)
	s.T().Log(url)
}

func (s *S3TestSuite) TestGetPresignedGetURL() {
	s.T().Skip()
	ctx := context.Background()

	url, err := s.cli.GeneratePresignedDownloadURL(ctx, "bagsy-notion", time.Minute*10)

	s.Require().NoError(err)
	s.T().Log(url)
}

func (s *S3TestSuite) TestGeneratePresignedPostURL() {
	ctx := context.Background()
	resp, err := s.cli.GeneratePresignedPostURL(ctx, UploadPolicyOptions{
		Key:              "bagsy-notion",
		ContentType:      "image/jpeg",
		ContentLengthMin: 1,
		ContentLengthMax: -1,
		Expires:          time.Minute * 15,
	})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.URL)
	s.T().Logf("resp: %v", resp)
}

// TestS3Suite — это стандартная Go test функция, которая запускает весь наш suite
func TestS3Suite(t *testing.T) {
	suite.Run(t, new(S3TestSuite))
}
