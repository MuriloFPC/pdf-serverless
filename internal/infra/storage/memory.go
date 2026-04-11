package storage

import (
	"context"
	"fmt"
	"sync"
)

type MemoryStorage struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (s *MemoryStorage) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[filename] = data
	return filename, nil
}

func (s *MemoryStorage) Download(ctx context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("file not found")
	}
	return data, nil
}

func (s *MemoryStorage) GetPresignedUploadURL(ctx context.Context, key string) (string, error) {
	return "http://localhost/upload/" + key, nil
}

func (s *MemoryStorage) GetPresignedDownloadURL(ctx context.Context, key string) (string, error) {
	return "http://localhost/download/" + key, nil
}

func (s *MemoryStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}
