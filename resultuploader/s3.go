package resultuploader

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/garnet-org/pkg/analysisrequest"
	"github.com/garnet-org/pkg/observability/tracer"
	"go.opentelemetry.io/otel/codes"
)

const S3UploaderName = "s3"

type S3 struct {
	region      string
	bucket      string
	endpoint    string
	credentials aws.CredentialsProvider
	s3Cfg       aws.Config
}

func NewS3Uploader(region string, bucket string, credentials aws.CredentialsProvider, endpoint string) *S3 {
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials,
	}
	if len(endpoint) > 0 {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})
	}
	return &S3{
		region:      region,
		bucket:      bucket,
		credentials: credentials,
		endpoint:    endpoint,
		s3Cfg:       cfg,
	}
}

func (s *S3) String() string {
	return S3UploaderName
}

func (s *S3) AlreadyExists(parent context.Context, uploadPath analysisrequest.ResultUploadPath) bool {
	ctx, span := tracer.FromContext(parent).Start(parent, "S3.AlreadyExists")
	defer span.End()

	client := s3.NewFromConfig(s.s3Cfg)

	fullKey := uploadPath.ToS3Key()

	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(fullKey),
	})
	return err == nil
}

func (s *S3) Upload(parent context.Context, src io.Reader, uploadPath analysisrequest.ResultUploadPath) error {
	ctx, span := tracer.FromContext(parent).Start(parent, "S3.Upload")
	defer span.End()

	s3Client := s3.NewFromConfig(s.s3Cfg)

	fullKey := uploadPath.ToS3Key()

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &fullKey,
		Body:        src,
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
