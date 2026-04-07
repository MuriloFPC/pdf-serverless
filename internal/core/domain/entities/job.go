package entities

import (
	"context"
	"time"
)

type JobStatus string

const (
	StatusAwaitingFiles         JobStatus = "awaiting_files"
	StatusPending               JobStatus = "pending"
	StatusProcessing            JobStatus = "processing"
	StatusCompleted             JobStatus = "completed"
	StatusFailed                JobStatus = "failed"
	StatusManuallyExcluded      JobStatus = "manually_excluded"
	StatusAutomaticallyExcluded JobStatus = "automatically_excluded"
)

type TTLType string

const (
	TTL24h     TTLType = "24h"
	TTL72h     TTLType = "72h"
	TTL1Week   TTLType = "1week"
	TTL1Month  TTLType = "1month"
	TTL3Month  TTLType = "3month"
	TTL6Month  TTLType = "6month"
	TTL1Year   TTLType = "1year"
	TTL3Year   TTLType = "3years"
	TTL5Year   TTLType = "5years"
	TTL10Year  TTLType = "10years"
	TTLForever TTLType = "unlimited"
)

type ProcessType string

const (
	TypeMerge          ProcessType = "merge"
	TypeSplit          ProcessType = "split"
	TypeProtect        ProcessType = "protect"
	TypeRemovePassword ProcessType = "remove_password"
	TypeOptimize       ProcessType = "optimize"
)

type FileMetadata struct {
	Path       string    `json:"path" dynamodbav:"path"`
	SizeKB     float64   `json:"size_kb" dynamodbav:"size_kb"`
	PageCount  int       `json:"page_count" dynamodbav:"page_count"`
	Filename   string    `json:"filename" dynamodbav:"filename"`
	UploadedAt time.Time `json:"uploaded_at" dynamodbav:"uploaded_at"`
}

type PDFJob struct {
	JobID       string         `json:"job_id" dynamodbav:"job_id"`
	UserID      string         `json:"user_id" dynamodbav:"user_id"`
	ProcessType ProcessType    `json:"process_type" dynamodbav:"process_type"`
	Status      JobStatus      `json:"status" dynamodbav:"status"`
	TTL         TTLType        `json:"ttl" omitzero:"" dynamodbav:"ttl"`
	CreatedAt   time.Time      `json:"created_at" dynamodbav:"created_at"`
	DeleteAt    time.Time      `json:"delete_at" omitzero:"" dynamodbav:"delete_at"`
	DeletedAt   time.Time      `json:"deleted_at" omitzero:"" dynamodbav:"deleted_at"`
	InputFiles  []FileMetadata `json:"input_files" dynamodbav:"input_files"`
	OutputFiles []FileMetadata `json:"output_files" dynamodbav:"output_files"`
	Password    string         `json:"password,omitempty" dynamodbav:"password,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty" dynamodbav:"metadata,omitempty"`
}

type PDFJobRepository interface {
	Create(ctx context.Context, job *PDFJob) error
	GetByID(ctx context.Context, id string) (*PDFJob, error)
	GetByUserID(ctx context.Context, userID string) ([]*PDFJob, error)
	Update(ctx context.Context, job *PDFJob) error
}
