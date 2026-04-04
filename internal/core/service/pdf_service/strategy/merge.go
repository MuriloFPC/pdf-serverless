package strategy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type MergeStrategy struct {
	storage interfaces.StorageProvider
}

func NewMergeStrategy(storage interfaces.StorageProvider) *MergeStrategy {
	return &MergeStrategy{storage: storage}
}

func (s *MergeStrategy) Process(ctx context.Context, job *entities.PDFJob) error {
	if len(job.InputFiles) == 0 {
		return fmt.Errorf("no input files to merge")
	}

	var readers []io.ReadSeeker
	for i := range job.InputFiles {
		input := &job.InputFiles[i]
		data, err := s.storage.Download(ctx, input.Path)
		if err != nil {
			return fmt.Errorf("failed to download input file %s: %w", input.Path, err)
		}
		UpdateInputMetadata(input, data, nil)
		readers = append(readers, bytes.NewReader(data))
	}

	var resultBuf bytes.Buffer
	// api.MergeRaw(rs []io.ReadSeeker, w io.Writer, conf *model.Configuration) is the signature for in-memory merging.
	// We pass nil for the configuration for now.
	if err := api.MergeRaw(readers, &resultBuf, false, nil); err != nil {
		return fmt.Errorf("failed to merge PDFs: %w", err)
	}

	resData := resultBuf.Bytes()
	outputKey := fmt.Sprintf("ttl/%s/%s/output/merged_%s.pdf", job.TTL, job.JobID, uuid.New().String())
	finalKey, err := s.storage.Upload(ctx, outputKey, resData)
	if err != nil {
		return fmt.Errorf("failed to upload merged PDF: %w", err)
	}

	job.OutputFiles = []entities.FileMetadata{
		NewFileMetadata(finalKey, fmt.Sprintf("merged_%s.pdf", job.JobID), resData, nil),
	}
	return nil
}

func (s *MergeStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeMerge
}
