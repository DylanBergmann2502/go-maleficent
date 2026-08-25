// internal/pkg/storage/s3.go
package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Storage provides object storage backed by Amazon S3 or an S3-compatible service.
type S3Storage struct {
	client   *s3.Client
	bucket   string
	location string
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

	return &S3Storage{
		client:   client,
		bucket:   storageConfig.Bucket,
		location: strings.Trim(storageConfig.Location, "/"),
	}, nil
}

// Save uploads an object and returns its storage key.
func (s *S3Storage) Save(ctx context.Context, filename, contentType string, body io.Reader) (string, error) {
	objectID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("generate storage object ID: %w", err)
	}

	filename = cleanFilename(filename)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	key := path.Join(s.location, objectID.String(), filename)
	if _, err := s.Put(ctx, key, body, contentType); err != nil {
		return "", err
	}

	return key, nil
}

// Put stores an object in the configured bucket.
func (s *S3Storage) Put(ctx context.Context, key string, body io.Reader, contentType string) (*s3.PutObjectOutput, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	return s.client.PutObject(ctx, input)
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

// URL returns a temporary URL for reading an object.
func (s *S3Storage) URL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	presigned, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = expiry
	})
	if err != nil {
		return "", err
	}

	return presigned.URL, nil
}

func cleanFilename(filename string) string {
	filename = strings.TrimSpace(strings.ReplaceAll(filename, "\\", "/"))
	filename = path.Base(filename)
	if filename == "." || filename == "/" || filename == "" {
		return "file"
	}
	return filename
}
