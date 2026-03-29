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
	Metadata map[string]any       `json:"metadata"`
	Files    []string             `json:"files"`
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

	job := &entities.PDFJob{
		JobID:       uuid.New().String(),
		UserID:      userID,
		ProcessType: req.Type,
		Status:      entities.StatusAwaitingFiles,
		CreatedAt:   time.Now(),
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

	key := fmt.Sprintf("inputs/%s/%s", uuid.New().String(), filename)

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
