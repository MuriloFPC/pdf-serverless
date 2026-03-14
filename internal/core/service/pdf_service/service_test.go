package pdf_service_test

import (
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/service/pdf_service"
	"pdf_serverless/internal/core/service/pdf_service/strategy"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/queue"
	"pdf_serverless/internal/infra/storage"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createTestPDF() []byte {
	return []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\ntrailer\n<< /Root 1 0 R >>\n%%EOF")
}

func TestPDFService_ProcessJob(t *testing.T) {
	ctx := t.Context()
	jobRepo := database.NewJobMemoryRepository()
	store := storage.NewMemoryStorage()
	q := queue.NewMemoryQueue()

	strategies := []strategy.ProcessingStrategy{
		strategy.NewMergeStrategy(store),
	}
	service := pdf_service.NewPDFService(jobRepo, store, q, strategies)

	job := &entities.PDFJob{
		JobID:       uuid.New().String(),
		UserID:      "test-user",
		ProcessType: entities.TypeMerge,
		Status:      entities.StatusPending,
		CreatedAt:   time.Now(),
		InputFiles:  []string{"test1.pdf", "test2.pdf"},
	}

	// Create dummy PDF files in storage
	dummyPDF := []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\ntrailer\n<< /Root 1 0 R >>\n%%EOF")
	store.Upload(ctx, "test1.pdf", dummyPDF)
	store.Upload(ctx, "test2.pdf", dummyPDF)

	err := jobRepo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	err = service.ProcessJob(ctx, job.JobID)
	if err != nil {
		t.Fatalf("Failed to process job: %v", err)
	}

	updatedJob, err := jobRepo.GetByID(ctx, job.JobID)
	if err != nil {
		t.Fatalf("Failed to get job: %v", err)
	}

	if updatedJob.Status != entities.StatusCompleted {
		t.Errorf("Expected status %s, got %s", entities.StatusCompleted, updatedJob.Status)
	}
}
