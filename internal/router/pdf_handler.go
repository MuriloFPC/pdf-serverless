package router

import (
	"fmt"
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
	Type     entities.ProcessType `form:"type"`
	Metadata map[string]any       `form:"metadata"`
}

func (h *PDFHandler) Process(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}

	processType := entities.ProcessType(c.FormValue("type"))
	if processType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Process type is required"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No files uploaded"})
	}

	var inputFiles []string
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
		}
		defer f.Close()

		data := make([]byte, file.Size)
		if _, err := f.Read(data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file"})
		}

		key := fmt.Sprintf("inputs/%s/%s", uuid.New().String(), file.Filename)
		s3Key, err := h.storage.Upload(c.Context(), key, data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to storage"})
		}
		inputFiles = append(inputFiles, s3Key)
	}

	job := &entities.PDFJob{
		JobID:       uuid.New().String(),
		UserID:      userID,
		ProcessType: processType,
		Status:      entities.StatusPending,
		CreatedAt:   time.Now(),
		InputFiles:  inputFiles,
	}

	if err := h.service.CreateJob(c.Context(), job); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create job"})
	}

	return c.Status(fiber.StatusAccepted).JSON(job)
}

func (h *PDFHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	jobID := c.Params("id")

	job, err := h.service.GetJobStatus(c.Context(), jobID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "PDFJob not found"})
	}

	if job.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	return c.JSON(job)
}

func (h *PDFHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	jobs, err := h.service.ListUserJobs(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list jobs"})
	}

	return c.JSON(jobs)
}
