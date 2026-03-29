package entities

import (
	"context"
	"time"
)

type JobStatus string

const (
	StatusAwaitingFiles JobStatus = "awaiting_files"
	StatusPending       JobStatus = "pending"
	StatusProcessing    JobStatus = "processing"
	StatusCompleted     JobStatus = "completed"
	StatusFailed        JobStatus = "failed"
)

type ProcessType string

const (
	TypeMerge          ProcessType = "merge"
	TypeSplit          ProcessType = "split"
	TypeProtect        ProcessType = "protect"
	TypeRemovePassword ProcessType = "remove_password"
)

type PDFJob struct {
	JobID       string         `json:"job_id" dynamodbav:"job_id"`
	UserID      string         `json:"user_id" dynamodbav:"user_id"`
	ProcessType ProcessType    `json:"process_type" dynamodbav:"process_type"`
	Status      JobStatus      `json:"status" dynamodbav:"status"`
	CreatedAt   time.Time      `json:"created_at" dynamodbav:"created_at"`
	DeleteAt    time.Time      `json:"delete_at" omitzero:"" dynamodbav:"delete_at"`
	InputFiles  []string       `json:"input_files" dynamodbav:"input_files"`
	OutputFiles []string       `json:"output_files" dynamodbav:"output_files"`
	Metadata    map[string]any `json:"metadata,omitempty" dynamodbav:"metadata,omitempty"`
}

type PDFJobRepository interface {
	Create(ctx context.Context, job *PDFJob) error
	GetByID(ctx context.Context, id string) (*PDFJob, error)
	GetByUserID(ctx context.Context, userID string) ([]*PDFJob, error)
	Update(ctx context.Context, job *PDFJob) error
}
