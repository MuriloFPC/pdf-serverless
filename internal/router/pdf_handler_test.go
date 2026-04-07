package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/service/pdf_service"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/queue"
	"pdf_serverless/internal/infra/storage"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestPDFHandler_GetDownloadURL_Errors(t *testing.T) {
	app := fiber.New()
	jobRepo := database.NewJobMemoryRepository()
	store := storage.NewMemoryStorage()
	q := queue.NewMemoryQueue()
	service := pdf_service.NewPDFService(jobRepo, store, q, nil)
	handler := NewPDFHandler(service, store)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user")
		return c.Next()
	})
	app.Get("/pdf/:id/download", handler.GetDownloadURL)

	tests := []struct {
		name           string
		status         entities.JobStatus
		deleteAt       time.Time
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Manually Excluded",
			status:         entities.StatusManuallyExcluded,
			expectedStatus: http.StatusNotFound,
			expectedError:  "File was manually deleted",
		},
		{
			name:           "Automatically Excluded",
			status:         entities.StatusAutomaticallyExcluded,
			expectedStatus: http.StatusNotFound,
			expectedError:  "File has expired and was automatically deleted",
		},
		{
			name:           "Expired Status Completed",
			status:         entities.StatusCompleted,
			deleteAt:       time.Now().Add(-1 * time.Hour),
			expectedStatus: http.StatusNotFound,
			expectedError:  "File has expired and was automatically deleted",
		},
		{
			name:           "Not Completed",
			status:         entities.StatusPending,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Job is not completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobID := "job-" + uuid.New().String()
			job := &entities.PDFJob{
				JobID:     jobID,
				UserID:    "test-user",
				Status:    tt.status,
				DeleteAt:  tt.deleteAt,
				CreatedAt: time.Now(),
			}
			jobRepo.Create(t.Context(), job)

			req := httptest.NewRequest("GET", "/pdf/"+jobID+"/download?filename=test.pdf", nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var body map[string]string
			json.NewDecoder(resp.Body).Decode(&body)
			if body["error"] != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, body["error"])
			}
		})
	}
}

func TestPDFHandler_GetDownloadURL_Success(t *testing.T) {
	app := fiber.New()
	jobRepo := database.NewJobMemoryRepository()
	store := storage.NewMemoryStorage()
	q := queue.NewMemoryQueue()
	service := pdf_service.NewPDFService(jobRepo, store, q, nil)
	handler := NewPDFHandler(service, store)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "test-user")
		return c.Next()
	})
	app.Get("/pdf/:id/download", handler.GetDownloadURL)

	jobID := "job-success"
	job := &entities.PDFJob{
		JobID:     jobID,
		UserID:    "test-user",
		Status:    entities.StatusCompleted,
		DeleteAt:  time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		OutputFiles: []entities.FileMetadata{
			{Path: "output.pdf", Filename: "output.pdf"},
		},
	}
	jobRepo.Create(t.Context(), job)

	req := httptest.NewRequest("GET", "/pdf/"+jobID+"/download?filename=output.pdf", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}
}
