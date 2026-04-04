package pdf_service

import (
	"context"
	"errors"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"
	"pdf_serverless/internal/core/service/pdf_service/strategy"
	"time"
)

type PDFService struct {
	jobRepo    entities.PDFJobRepository
	storage    interfaces.StorageProvider
	queue      interfaces.QueueProvider
	strategies []strategy.ProcessingStrategy
}

func NewPDFService(
	jobRepo entities.PDFJobRepository,
	storage interfaces.StorageProvider,
	queue interfaces.QueueProvider,
	strategies []strategy.ProcessingStrategy,
) *PDFService {
	return &PDFService{
		jobRepo:    jobRepo,
		storage:    storage,
		queue:      queue,
		strategies: strategies,
	}
}

func (s *PDFService) CreateJob(ctx context.Context, job *entities.PDFJob) error {
	return s.jobRepo.Create(ctx, job)
}

func (s *PDFService) PublishJob(ctx context.Context, jobID string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	job.Status = entities.StatusPending
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return err
	}

	return s.queue.Publish(ctx, job.JobID)
}

func (s *PDFService) AddInputFile(ctx context.Context, jobID string, key string, filename string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	job.InputFiles = append(job.InputFiles, entities.FileMetadata{
		Path:       key,
		Filename:   filename,
		UploadedAt: time.Now(),
	})
	return s.jobRepo.Update(ctx, job)
}

func (s *PDFService) GetJobStatus(ctx context.Context, jobID string) (*entities.PDFJob, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

func (s *PDFService) ListUserJobs(ctx context.Context, userID string) ([]*entities.PDFJob, error) {
	return s.jobRepo.GetByUserID(ctx, userID)
}

func (s *PDFService) ProcessJob(ctx context.Context, jobID string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	var activeStrategy strategy.ProcessingStrategy
	for _, st := range s.strategies {
		if st.CanHandle(job.ProcessType) {
			activeStrategy = st
			break
		}
	}

	if activeStrategy == nil {
		return errors.New("no strategy found for process type")
	}

	job.Status = entities.StatusProcessing
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return err
	}

	if err := activeStrategy.Process(ctx, job); err != nil {
		job.Status = entities.StatusFailed
		job.Password = "" // Clear sensitive info
		s.jobRepo.Update(ctx, job)
		return err
	}

	job.Status = entities.StatusCompleted
	job.Password = "" // Clear sensitive info
	return s.jobRepo.Update(ctx, job)
}
