package interfaces

import (
	"context"
)

type StorageProvider interface {
	Upload(ctx context.Context, filename string, data []byte) (string, error)
	Download(ctx context.Context, key string) ([]byte, error)
	GetPresignedUploadURL(ctx context.Context, key string) (string, error)
	GetPresignedDownloadURL(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
