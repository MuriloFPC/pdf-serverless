package pdf_service

import (
	"context"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/storage"
	"testing"
)

func TestDeleteJobFiles(t *testing.T) {
	ctx := context.Background()
	jobRepo := database.NewJobMemoryRepository()
	store := storage.NewMemoryStorage()
	service := NewPDFService(jobRepo, store, nil, nil)

	userID := "user123"
	jobID := "job456"

	// Create a mock job
	job := &entities.PDFJob{
		JobID:  jobID,
		UserID: userID,
		InputFiles: []entities.FileMetadata{
			{Path: "input/test1.pdf", Filename: "test1.pdf"},
		},
		OutputFiles: []entities.FileMetadata{
			{Path: "output/result.pdf", Filename: "result.pdf"},
		},
		Status: entities.StatusCompleted,
	}

	err := jobRepo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	// Upload mock files to memory storage
	_, _ = store.Upload(ctx, "input/test1.pdf", []byte("pdf data"))
	_, _ = store.Upload(ctx, "output/result.pdf", []byte("result data"))

	// Verify files exist before deletion
	_, err = store.Download(ctx, "input/test1.pdf")
	if err != nil {
		t.Errorf("Input file should exist before deletion")
	}
	_, err = store.Download(ctx, "output/result.pdf")
	if err != nil {
		t.Errorf("Output file should exist before deletion")
	}

	// Delete job files
	err = service.DeleteJobFiles(ctx, userID, jobID)
	if err != nil {
		t.Fatalf("DeleteJobFiles failed: %v", err)
	}

	// Verify job status is updated
	updatedJob, _ := jobRepo.GetByID(ctx, jobID)
	if updatedJob.Status != entities.StatusManuallyExcluded {
		t.Errorf("Expected status %s, got %s", entities.StatusManuallyExcluded, updatedJob.Status)
	}

	// Verify files are deleted from storage
	_, err = store.Download(ctx, "input/test1.pdf")
	if err == nil {
		t.Errorf("Input file should be deleted from storage")
	}

	_, err = store.Download(ctx, "output/result.pdf")
	if err == nil {
		t.Errorf("Output file should be deleted from storage")
	}
}
