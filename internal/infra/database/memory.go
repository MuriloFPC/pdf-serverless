package database

import (
	"context"
	"errors"
	"pdf_serverless/internal/core/domain/entities"
	"sync"
)

type JobMemoryRepository struct {
	jobs map[string]*entities.PDFJob
	mu   sync.RWMutex
}

func NewJobMemoryRepository() *JobMemoryRepository {
	return &JobMemoryRepository{
		jobs: make(map[string]*entities.PDFJob),
	}
}

func (r *JobMemoryRepository) Create(ctx context.Context, job *entities.PDFJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.JobID] = job
	return nil
}

func (r *JobMemoryRepository) GetByID(ctx context.Context, id string) (*entities.PDFJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[id]
	if !ok {
		return nil, errors.New("job not found")
	}
	return job, nil
}

func (r *JobMemoryRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.PDFJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var jobs []*entities.PDFJob
	for _, job := range r.jobs {
		if job.UserID == userID {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func (r *JobMemoryRepository) Update(ctx context.Context, job *entities.PDFJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.JobID] = job
	return nil
}

type UserMemoryRepository struct {
	users map[string]*entities.User
	mu    sync.RWMutex
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users: make(map[string]*entities.User),
	}
}

func (r *UserMemoryRepository) Create(ctx context.Context, user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *UserMemoryRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserMemoryRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
