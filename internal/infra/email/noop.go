package email

import (
	"context"
	"pdf_serverless/internal/core/domain/interfaces"
)

type NoOpEmailService struct{}

func NewNoOpEmailService() interfaces.EmailService {
	return &NoOpEmailService{}
}

func (s *NoOpEmailService) Send(ctx context.Context, to, subject, body string) error {
	return nil
}
