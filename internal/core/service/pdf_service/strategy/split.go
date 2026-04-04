package strategy

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type SplitStrategy struct {
	storage interfaces.StorageProvider
}

func NewSplitStrategy(storage interfaces.StorageProvider) *SplitStrategy {
	return &SplitStrategy{storage: storage}
}

func (s *SplitStrategy) Process(ctx context.Context, job *entities.PDFJob) error {
	if len(job.InputFiles) == 0 {
		return fmt.Errorf("no input files to split")
	}

	input := &job.InputFiles[0]
	data, err := s.storage.Download(ctx, input.Path)
	if err != nil {
		return fmt.Errorf("failed to download input file %s: %w", input.Path, err)
	}

	UpdateInputMetadata(input, data, nil)
	rs := bytes.NewReader(data)

	// api.SplitRaw(rs io.ReadSeeker, outDir string, fileName string, conf *model.Configuration)
	// Actually SplitRaw still takes an outDir.

	// Let's see if there is an alternative in pdfcpu to get buffers back.
	// Looking at pdfcpu source, Split usually writes to disk.

	// For MVP, we can use a temporary directory.
	tempDir, err := os.MkdirTemp("", "pdf-split-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if err := api.Split(rs, tempDir, "split", 1, nil); err != nil {
		return fmt.Errorf("failed to split PDF: %w", err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read split files: %w", err)
	}

	var outputFiles []entities.FileMetadata
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(tempDir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read split file %s: %w", file.Name(), err)
		}

		outputKey := fmt.Sprintf("ttl/%s/%s/output/%s", job.TTL, job.JobID, file.Name())
		finalKey, err := s.storage.Upload(ctx, outputKey, data)
		if err != nil {
			return fmt.Errorf("failed to upload split file %s: %w", file.Name(), err)
		}

		outputFiles = append(outputFiles, NewFileMetadata(finalKey, file.Name(), data, nil))
	}

	// Zip files if there are more than one
	if len(outputFiles) > 1 {
		zipBuffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(zipBuffer)

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			path := filepath.Join(tempDir, file.Name())
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file for zipping %s: %w", file.Name(), err)
			}

			w, err := zipWriter.Create(file.Name())
			if err != nil {
				f.Close()
				return fmt.Errorf("failed to create zip entry for %s: %w", file.Name(), err)
			}

			if _, err := io.Copy(w, f); err != nil {
				f.Close()
				return fmt.Errorf("failed to copy file to zip %s: %w", file.Name(), err)
			}
			f.Close()
		}

		if err := zipWriter.Close(); err != nil {
			return fmt.Errorf("failed to close zip writer: %w", err)
		}

		zipData := zipBuffer.Bytes()
		zipKey := fmt.Sprintf("ttl/%s/%s/output/split_results.zip", job.TTL, job.JobID)
		finalZipKey, err := s.storage.Upload(ctx, zipKey, zipData)
		if err != nil {
			return fmt.Errorf("failed to upload zip file: %w", err)
		}

		outputFiles = append(outputFiles, NewFileMetadata(finalZipKey, "split_results.zip", zipData, nil))
	}

	job.OutputFiles = outputFiles
	return nil
}

func (s *SplitStrategy) CanHandle(processType entities.ProcessType) bool {
	return processType == entities.TypeSplit
}
