package interfaces

import (
	"context"
)

type StorageProvider interface {
	Upload(ctx context.Context, filename string, data []byte) (string, error)
	Download(ctx context.Context, key string) ([]byte, error)
}
