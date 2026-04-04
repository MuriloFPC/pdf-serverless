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

	input := &job.InputFiles[0]
	data, err := s.storage.Download(ctx, input.Path)
	if err != nil {
		return fmt.Errorf("failed to download input file %s: %w", input.Path, err)
	}

	conf := model.NewDefaultConfiguration()
	conf.UserPW = job.Password
	conf.OwnerPW = job.Password

	UpdateInputMetadata(input, data, conf)
	rs := bytes.NewReader(data)

	var resultBuf bytes.Buffer

	if err := api.Decrypt(rs, &resultBuf, conf); err != nil {
		return fmt.Errorf("failed to unprotect PDF: %w", err)
	}

	resData := resultBuf.Bytes()
	outputKey := fmt.Sprintf("ttl/%s/%s/output/unprotected_%s.pdf", job.TTL, job.JobID, uuid.New().String())
	finalKey, err := s.storage.Upload(ctx, outputKey, resData)
	if err != nil {
		return fmt.Errorf("failed to upload unprotected PDF: %w", err)
	}

	job.OutputFiles = []entities.FileMetadata{
		NewFileMetadata(finalKey, fmt.Sprintf("unprotected_%s.pdf", job.JobID), resData, nil),
	}
	return nil
}

func (s *UnprotectStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeRemovePassword
}
