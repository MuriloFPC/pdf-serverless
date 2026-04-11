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

// Process godoc
// @Summary      Create a new PDF processing job
// @Description  Initialize a PDF job (merge, split, protect, etc.)
// @Tags         pdf
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string          true  "Bearer <token>"
// @Param        request        body      ProcessRequest  true  "Process request"
// @Success      201            {object}  entities.PDFJob
// @Failure      400            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /pdf/process [post]
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

	if req.TTL == "" || req.TTL == "6h" {
		req.TTL = entities.TTL24h // Default TTL
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

// GetPresignedURL godoc
// @Summary      Get a presigned S3 upload URL
// @Description  Generates a URL for direct upload to S3 for a specific job
// @Tags         pdf
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Param        id             path      string  true  "Job ID"
// @Param        filename       query     string  true  "Original filename"
// @Success      200            {object}  map[string]string
// @Failure      400            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Failure      403            {object}  map[string]string
// @Failure      404            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /pdf/presigned-url/{id} [get]
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

	if err := h.service.AddInputFile(c.Context(), jobID, key, filename); err != nil {
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

// CompleteUpload godoc
// @Summary      Mark job uploads as complete
// @Description  Notifies the system that all files for a job have been uploaded
// @Tags         pdf
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Param        id             path      string  true  "Job ID"
// @Success      200            {object}  map[string]string
// @Failure      400            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Failure      403            {object}  map[string]string
// @Failure      404            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /pdf/complete-upload/{id} [post]
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

// GetStatus godoc
// @Summary      Get PDF job status
// @Description  Returns the current status of a PDF processing job
// @Tags         pdf
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Param        id             path      string  true  "Job ID"
// @Success      200            {object}  entities.PDFJob
// @Failure      401            {object}  map[string]string
// @Failure      403            {object}  map[string]string
// @Failure      404            {object}  map[string]string
// @Router       /pdf/status/{id} [get]
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

	// Update status if it was completed but is now past expiration time
	if job.Status == entities.StatusCompleted && !job.DeleteAt.IsZero() && time.Now().After(job.DeleteAt) {
		job.Status = entities.StatusAutomaticallyExcluded
	}

	return c.JSON(job)
}

// List godoc
// @Summary      List user PDF jobs
// @Description  Returns all PDF processing jobs for the authenticated user
// @Tags         pdf
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Success      200            {array}   entities.PDFJob
// @Failure      401            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /pdf/list [get]
func (h *PDFHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	jobs, err := h.service.ListUserJobs(c.Context(), userID)
	if err != nil {
		log.Printf("PDFHandler.List: Error listing jobs for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list jobs"})
	}

	now := time.Now()
	for i := range jobs {
		if jobs[i].Status == entities.StatusCompleted && !jobs[i].DeleteAt.IsZero() && now.After(jobs[i].DeleteAt) {
			jobs[i].Status = entities.StatusAutomaticallyExcluded
		}
	}

	return c.JSON(jobs)
}

// GetDownloadURL godoc
// @Summary      Get a presigned S3 download URL
// @Description  Generates a URL for downloading a processed PDF file
// @Tags         pdf
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Param        id             path      string  true  "Job ID"
// @Param        filename       query     string  true  "File path from job output"
// @Success      200            {object}  map[string]string
// @Failure      400            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Failure      403            {object}  map[string]string
// @Failure      404            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /pdf/download/{id} [get]
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
		if job.Status == entities.StatusManuallyExcluded {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File was manually deleted"})
		}
		if job.Status == entities.StatusAutomaticallyExcluded || (!job.DeleteAt.IsZero() && time.Now().After(job.DeleteAt)) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File has expired and was automatically deleted"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job is not completed"})
	}

	// Update status if it was completed but is now past expiration time (double check)
	if !job.DeleteAt.IsZero() && time.Now().After(job.DeleteAt) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File has expired and was automatically deleted"})
	}

	// Find the file in output_files
	var found bool
	for _, f := range job.OutputFiles {
		if f.Path == filename {
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

	log.Printf("PDFHandler.GetDownloadURL: User %s is downloading file %s from job %s", userID, filename, jobID)

	return c.JSON(fiber.Map{
		"url": url,
	})
}

// Delete godoc
// @Summary      Delete a PDF job
// @Description  Removes a PDF job and its associated files
// @Tags         pdf
// @Param        Authorization  header    string  true  "Bearer <token>"
// @Param        id             path      string  true  "Job ID"
// @Success      204            "No Content"
// @Failure      401            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /pdf/{id} [delete]
func (h *PDFHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jobID := c.Params("id")

	if err := h.service.DeleteJobFiles(c.Context(), userID, jobID); err != nil {
		log.Printf("PDFHandler.Delete: Error deleting job %s: %v", jobID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete job files"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func calculateDeleteAt(createdAt time.Time, ttl entities.TTLType) time.Time {
	var duration time.Duration
	switch ttl {
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
		duration = 24 * time.Hour // Default to 24h if unknown
	}
	return createdAt.Add(duration)
}
