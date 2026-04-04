package strategy

import (
	"bytes"
	"context"
	"pdf_serverless/internal/core/domain/entities"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type ProcessingStrategy interface {
	Process(ctx context.Context, job *entities.PDFJob) error
	CanHandle(processType entities.ProcessType) bool
}

func NewFileMetadata(path string, filename string, data []byte, conf *model.Configuration) entities.FileMetadata {
	pCount := 0
	if data != nil {
		rs := bytes.NewReader(data)
		pCount, _ = api.PageCount(rs, conf)
	}

	return entities.FileMetadata{
		Path:       path,
		Filename:   filename,
		SizeKB:     float64(len(data)) / 1024.0,
		PageCount:  pCount,
		UploadedAt: time.Now(),
	}
}

func UpdateInputMetadata(input *entities.FileMetadata, data []byte, conf *model.Configuration) {
	if data == nil {
		return
	}
	rs := bytes.NewReader(data)
	pCount, _ := api.PageCount(rs, conf)
	input.PageCount = pCount
	input.SizeKB = float64(len(data)) / 1024.0
}
