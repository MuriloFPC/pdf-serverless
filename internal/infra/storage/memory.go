package storage

import (
	"context"
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
		return nil, nil
	}
	return data, nil
}

func (s *MemoryStorage) GetPresignedUploadURL(ctx context.Context, key string) (string, error) {
	return "http://localhost/upload/" + key, nil
}
