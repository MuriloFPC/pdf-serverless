package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(client *s3.Client) *S3Storage {
	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		bucket = "pdf-serverless-opensource"
	}
	return &S3Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *S3Storage) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return filename, nil
}

func (s *S3Storage) Download(ctx context.Context, key string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}
	defer out.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(out.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
	}

	return buf.Bytes(), nil
}
