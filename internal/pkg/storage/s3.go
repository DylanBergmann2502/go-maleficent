// internal/pkg/storage/s3.go
package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage provides object storage backed by Amazon S3 or an S3-compatible service.
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates an S3 client configured for Amazon S3, Garage, or another
// S3-compatible object storage service.
func NewS3Storage(ctx context.Context, storageConfig *config.StorageConfig) (*S3Storage, error) {
	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(storageConfig.Region),
	}
	if storageConfig.AccessKeyID != "" && storageConfig.SecretAccessKey != "" {
		loadOptions = append(loadOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				storageConfig.AccessKeyID,
				storageConfig.SecretAccessKey,
				"",
			),
		))
	}

	sdkConfig, err := awsconfig.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, fmt.Errorf("load S3 configuration: %w", err)
	}

	client := s3.NewFromConfig(sdkConfig, func(options *s3.Options) {
		options.UsePathStyle = storageConfig.UsePathStyle
		if storageConfig.Endpoint != "" {
			options.BaseEndpoint = aws.String(storageConfig.Endpoint)
		}
	})

	return &S3Storage{client: client, bucket: storageConfig.Bucket}, nil
}

// Put stores an object in the configured bucket.
func (s *S3Storage) Put(ctx context.Context, key string, body io.Reader, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err := s.client.PutObject(ctx, input)
	return err
}

// Get retrieves an object from the configured bucket. The caller must close the
// returned body.
func (s *S3Storage) Get(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	return s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
}

// Delete removes an object from the configured bucket.
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// Health checks whether the configured bucket can be reached.
func (s *S3Storage) Health(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	return err
}
