package strategy

import (
	"bytes"
	"context"
	"fmt"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type OptimizeStrategy struct {
	storage interfaces.StorageProvider
}

func NewOptimizeStrategy(storage interfaces.StorageProvider) *OptimizeStrategy {
	return &OptimizeStrategy{storage: storage}
}

func (s *OptimizeStrategy) Process(ctx context.Context, job *entities.PDFJob) error {
	if len(job.InputFiles) == 0 {
		return fmt.Errorf("no input files to optimize")
	}

	// Otimizamos o primeiro arquivo de entrada (ou poderíamos fazer para todos se for o caso,
	// mas geralmente otimização é por arquivo. Se houver mais de um,
	// assumimos que o usuário quer otimizar o primeiro ou o sistema foi desenhado para processar um por vez para este tipo.)
	input := &job.InputFiles[0]
	data, err := s.storage.Download(ctx, input.Path)
	if err != nil {
		return fmt.Errorf("failed to download input file %s: %w", input.Path, err)
	}
	UpdateInputMetadata(input, data, nil)

	var resultBuf bytes.Buffer
	rs := bytes.NewReader(data)

	if err := api.Optimize(rs, &resultBuf, nil); err != nil {
		return fmt.Errorf("failed to optimize PDF: %w", err)
	}

	resData := resultBuf.Bytes()
	outputKey := fmt.Sprintf("ttl/%s/%s/output/optimized_%s.pdf", job.TTL, job.JobID, uuid.New().String())
	finalKey, err := s.storage.Upload(ctx, outputKey, resData)
	if err != nil {
		return fmt.Errorf("failed to upload optimized PDF: %w", err)
	}

	job.OutputFiles = []entities.FileMetadata{
		NewFileMetadata(finalKey, fmt.Sprintf("optimized_%s.pdf", job.JobID), resData, nil),
	}

	return nil
}

func (s *OptimizeStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeOptimize
}
