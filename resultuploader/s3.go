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
	bucket string
	s3Cfg  aws.Config
}

func NewS3Uploader(cfg aws.Config, bucket string) *S3 {
	return &S3{
		bucket: bucket,
		s3Cfg:  cfg,
	}
}

func (s *S3) WithCustomEndpoint(endpoint string) *S3 {
	s.s3Cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               endpoint,
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})
	return s
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

// This is needed because the s3.GetObjectOutput does not retain any information about the key.
type S3GetObject struct {
	*s3.GetObjectOutput
	Key string
}

func (s *S3) GetObject(parent context.Context, key string) (*S3GetObject, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "S3.GetObject")
	defer span.End()

	s3Client := s3.NewFromConfig(s.s3Cfg)

	obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &S3GetObject{obj, key}, nil
}
