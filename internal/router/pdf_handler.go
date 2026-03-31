package router

import (
	"fmt"
	"log"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"
	"pdf_serverless/internal/core/service/pdf_service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PDFHandler struct {
	service *pdf_service.PDFService
	storage interfaces.StorageProvider
}

func NewPDFHandler(service *pdf_service.PDFService, storage interfaces.StorageProvider) *PDFHandler {
	return &PDFHandler{
		service: service,
		storage: storage,
	}
}

type ProcessRequest struct {
	Type     entities.ProcessType `json:"type"`
	TTL      entities.TTLType     `json:"ttl"`
	Password string               `json:"password,omitempty"`
	Metadata map[string]any       `json:"metadata"`
}

func (h *PDFHandler) Process(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req ProcessRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("PDFHandler.Process: Error parsing body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse body"})
	}

	if req.Type == "" {
		log.Printf("PDFHandler.Process: Missing process type")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Process type is required"})
	}

	if (req.Type == entities.TypeProtect || req.Type == entities.TypeRemovePassword) && req.Password == "" {
		log.Printf("PDFHandler.Process: Password required for %s", req.Type)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password is required for this process type"})
	}

	if req.TTL == "" {
		req.TTL = entities.TTL6h // Default TTL
	}

	createdAt := time.Now()
	deleteAt := calculateDeleteAt(createdAt, req.TTL)

	job := &entities.PDFJob{
		JobID:       uuid.New().String(),
		UserID:      userID,
		ProcessType: req.Type,
		Status:      entities.StatusAwaitingFiles,
		TTL:         req.TTL,
		CreatedAt:   createdAt,
		DeleteAt:    deleteAt,
		Password:    req.Password,
		Metadata:    req.Metadata,
	}

	if err := h.service.CreateJob(c.Context(), job); err != nil {
		log.Printf("PDFHandler.Process: Error creating job in service: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create job"})
	}

	return c.Status(fiber.StatusCreated).JSON(job)
}

func (h *PDFHandler) GetPresignedURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jobID := c.Params("id")
	filename := c.Query("filename")

	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filename is required"})
	}

	job, err := h.service.GetJobStatus(c.Context(), jobID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
	}

	if job.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	if job.Status != entities.StatusAwaitingFiles {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job is not awaiting files"})
	}

	key := fmt.Sprintf("ttl/%s/%s/input/%s_%s.pdf", job.TTL, jobID, filename, uuid.New().String())

	if err := h.service.AddInputFile(c.Context(), jobID, key); err != nil {
		log.Printf("PDFHandler.GetPresignedURL: Error adding input file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update job with new file"})
	}

	url, err := h.storage.GetPresignedUploadURL(c.Context(), key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presigned URL"})
	}

	return c.JSON(fiber.Map{
		"url": url,
		"key": key,
	})
}

func (h *PDFHandler) CompleteUpload(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jobID := c.Params("id")

	job, err := h.service.GetJobStatus(c.Context(), jobID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
	}

	if job.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	if job.Status != entities.StatusAwaitingFiles {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job is not awaiting files"})
	}

	if err := h.service.PublishJob(c.Context(), jobID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to complete upload"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "pending"})
}

func (h *PDFHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jobID := c.Params("id")

	job, err := h.service.GetJobStatus(c.Context(), jobID)
	if err != nil {
		log.Printf("PDFHandler.GetStatus: Job %s not found: %v", jobID, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "PDFJob not found"})
	}

	if job.UserID != userID {
		log.Printf("PDFHandler.GetStatus: Access denied for user %s to job %s", userID, jobID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	return c.JSON(job)
}

func (h *PDFHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	jobs, err := h.service.ListUserJobs(c.Context(), userID)
	if err != nil {
		log.Printf("PDFHandler.List: Error listing jobs for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list jobs"})
	}

	return c.JSON(jobs)
}

func (h *PDFHandler) GetDownloadURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jobID := c.Params("id")
	filename := c.Query("filename")

	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Filename is required"})
	}

	job, err := h.service.GetJobStatus(c.Context(), jobID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
	}

	if job.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	if job.Status != entities.StatusCompleted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job is not completed"})
	}

	// Find the file in output_files
	var found bool
	for _, f := range job.OutputFiles {
		if f == filename {
			found = true
			break
		}
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found in job outputs"})
	}

	url, err := h.storage.GetPresignedDownloadURL(c.Context(), filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate download URL"})
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}

func calculateDeleteAt(createdAt time.Time, ttl entities.TTLType) time.Time {
	var duration time.Duration
	switch ttl {
	case entities.TTL6h:
		duration = 6 * time.Hour
	case entities.TTL24h:
		duration = 24 * time.Hour
	case entities.TTL72h:
		duration = 72 * time.Hour
	case entities.TTL1Week:
		duration = 7 * 24 * time.Hour
	case entities.TTL1Month:
		duration = 30 * 24 * time.Hour
	case entities.TTL3Month:
		duration = 90 * 24 * time.Hour
	case entities.TTL6Month:
		duration = 180 * 24 * time.Hour
	case entities.TTL1Year:
		duration = 365 * 24 * time.Hour
	case entities.TTL3Year:
		duration = 3 * 365 * 24 * time.Hour
	case entities.TTL5Year:
		duration = 5 * 365 * 24 * time.Hour
	case entities.TTL10Year:
		duration = 10 * 365 * 24 * time.Hour
	case entities.TTLForever:
		return time.Time{} // No expiration
	default:
		duration = 6 * time.Hour // Default to 6h if unknown
	}
	return createdAt.Add(duration)
}
