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

type ProtectStrategy struct {
	storage interfaces.StorageProvider
}

func NewProtectStrategy(storage interfaces.StorageProvider) *ProtectStrategy {
	return &ProtectStrategy{storage: storage}
}

func (s *ProtectStrategy) Process(ctx context.Context, job *entities.PDFJob) error {
	if len(job.InputFiles) == 0 {
		return fmt.Errorf("no input files to protect")
	}

	if job.Password == "" {
		return fmt.Errorf("password is required for protect operation")
	}

	input := &job.InputFiles[0]
	data, err := s.storage.Download(ctx, input.Path)
	if err != nil {
		return fmt.Errorf("failed to download input file %s: %w", input.Path, err)
	}

	UpdateInputMetadata(input, data, nil)
	rs := bytes.NewReader(data)

	var resultBuf bytes.Buffer

	// Configure encryption
	conf := model.NewDefaultConfiguration()
	conf.UserPW = job.Password
	conf.OwnerPW = job.Password
	// You can also set other permissions here if needed

	if err := api.Encrypt(rs, &resultBuf, conf); err != nil {
		return fmt.Errorf("failed to protect PDF: %w", err)
	}

	resData := resultBuf.Bytes()
	outputKey := fmt.Sprintf("ttl/%s/%s/output/protected_%s.pdf", job.TTL, job.JobID, uuid.New().String())
	finalKey, err := s.storage.Upload(ctx, outputKey, resData)
	if err != nil {
		return fmt.Errorf("failed to upload protected PDF: %w", err)
	}

	job.OutputFiles = []entities.FileMetadata{
		NewFileMetadata(finalKey, fmt.Sprintf("protected_%s.pdf", job.JobID), resData, conf),
	}
	return nil
}

func (s *ProtectStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeProtect
}
