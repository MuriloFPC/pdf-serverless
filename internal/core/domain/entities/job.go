package entities

import (
	"context"
	"time"
)

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type ProcessType string

const (
	TypeMerge          ProcessType = "merge"
	TypeSplit          ProcessType = "split"
	TypeProtect        ProcessType = "protect"
	TypeRemovePassword ProcessType = "remove_password"
)

type PDFJob struct {
	JobID       string         `json:"job_id"`
	UserID      string         `json:"user_id"`
	ProcessType ProcessType    `json:"process_type"`
	Status      JobStatus      `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	DeleteAt    time.Time      `json:"delete_at" omitzero:""`
	InputFiles  []string       `json:"input_files"`
	OutputFiles []string       `json:"output_files"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type PDFJobRepository interface {
	Create(ctx context.Context, job *PDFJob) error
	GetByID(ctx context.Context, id string) (*PDFJob, error)
	GetByUserID(ctx context.Context, userID string) ([]*PDFJob, error)
	Update(ctx context.Context, job *PDFJob) error
}
