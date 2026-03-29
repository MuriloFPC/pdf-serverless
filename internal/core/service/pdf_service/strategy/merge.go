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
	for _, key := range job.InputFiles {
		data, err := s.storage.Download(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to download input file %s: %w", key, err)
		}
		readers = append(readers, bytes.NewReader(data))
	}

	var resultBuf bytes.Buffer
	// api.MergeRaw(rs []io.ReadSeeker, w io.Writer, conf *model.Configuration) is the signature for in-memory merging.
	// We pass nil for the configuration for now.
	if err := api.MergeRaw(readers, &resultBuf, false, nil); err != nil {
		return fmt.Errorf("failed to merge PDFs: %w", err)
	}

	outputKey := fmt.Sprintf("%s/output/merged_%s.pdf", job.JobID, uuid.New().String())
	finalKey, err := s.storage.Upload(ctx, outputKey, resultBuf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to upload merged PDF: %w", err)
	}

	job.OutputFiles = []string{finalKey}
	return nil
}

func (s *MergeStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeMerge
}
