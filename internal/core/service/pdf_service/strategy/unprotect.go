package strategy

import (
	"bytes"
	"context"
	"fmt"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type UnprotectStrategy struct {
	storage interfaces.StorageProvider
}

func NewUnprotectStrategy(storage interfaces.StorageProvider) *UnprotectStrategy {
	return &UnprotectStrategy{storage: storage}
}

func (s *UnprotectStrategy) Process(ctx context.Context, job *entities.PDFJob) error {
	if len(job.InputFiles) == 0 {
		return fmt.Errorf("no input files to unprotect")
	}

	key := job.InputFiles[0]
	data, err := s.storage.Download(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to download input file %s: %w", key, err)
	}

	rs := bytes.NewReader(data)
	var resultBuf bytes.Buffer

	conf := model.NewDefaultConfiguration()
	conf.UserPW = job.Password
	conf.OwnerPW = job.Password

	if err := api.Decrypt(rs, &resultBuf, conf); err != nil {
		return fmt.Errorf("failed to unprotect PDF: %w", err)
	}

	outputKey := fmt.Sprintf("ttl/%s/%s/output/unprotected_%s.pdf", job.TTL, job.JobID, uuid.New().String())
	finalKey, err := s.storage.Upload(ctx, outputKey, resultBuf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to upload unprotected PDF: %w", err)
	}

	job.OutputFiles = []string{finalKey}
	return nil
}

func (s *UnprotectStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeRemovePassword
}
