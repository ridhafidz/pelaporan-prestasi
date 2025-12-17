package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/repository"
)

// StudentService defines business logic operations for students (routes expect Find*/Update with id)
type StudentService interface {
	FindAll(ctx context.Context) ([]models.Student, error)
	FindByID(ctx context.Context, id string) (*models.Student, error)
	Create(ctx context.Context, s *models.Student) error
	Update(ctx context.Context, id string, s *models.Student) error
	Delete(ctx context.Context, id string) error
}

type studentService struct {
	repo repository.StudentRepository
}

// NewStudentService creates a new StudentService
func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (svc *studentService) FindAll(ctx context.Context) ([]models.Student, error) {
	return svc.repo.GetAll(ctx)
}

func (svc *studentService) FindByID(ctx context.Context, id string) (*models.Student, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *studentService) Create(ctx context.Context, s *models.Student) error {
	if s == nil {
		return errors.New("student is nil")
	}
	if s.Name == "" {
		return errors.New("name is required")
	}
	if s.Email == "" {
		return errors.New("email is required")
	}

	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	s.UpdatedAt = now

	// ensure ID if missing
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	return svc.repo.Create(ctx, s)
}

func (svc *studentService) Update(ctx context.Context, id string, s *models.Student) error {
	if s == nil {
		return errors.New("student is nil")
	}
	if id == "" {
		return errors.New("student id is required")
	}

	// try parse id to uuid and set on model if possible (repository expects id param or uses model's id)
	if parsed, err := uuid.Parse(id); err == nil {
		s.ID = parsed
	}

	s.UpdatedAt = time.Now()
	return svc.repo.Update(ctx, s)
}

func (svc *studentService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return svc.repo.Delete(ctx, id)
}
