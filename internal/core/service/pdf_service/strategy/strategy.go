package strategy

import (
	"context"
	"pdf_serverless/internal/core/domain/entities"
)

type ProcessingStrategy interface {
	Process(ctx context.Context, job *entities.PDFJob) error
	CanHandle(processType entities.ProcessType) bool
}
