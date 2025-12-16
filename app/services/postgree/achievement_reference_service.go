package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/repository"
)

// AchievementReferenceService contains business logic for achievement references
type AchievementReferenceService interface {
	Create(ctx context.Context, req *models.CreateAchievementReferenceRequest) (*models.AchievementReference, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.AchievementStatus, verifiedBy *uuid.UUID, rejectionNote *string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type achievementReferenceService struct {
	repo repository.AchievementReferenceRepository
}

// NewAchievementReferenceService creates a new service instance
func NewAchievementReferenceService(repo repository.AchievementReferenceRepository) AchievementReferenceService {
	return &achievementReferenceService{repo: repo}
}

func (s *achievementReferenceService) Create(ctx context.Context, req *models.CreateAchievementReferenceRequest) (*models.AchievementReference, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.StudentID == uuid.Nil {
		return nil, errors.New("studentId is required")
	}
	if req.MongoAchievementID == "" {
		return nil, errors.New("mongoAchievementId is required")
	}

	now := time.Now()
	ref := &models.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          req.StudentID,
		MongoAchievementID: req.MongoAchievementID,
		Status:             models.StatusSubmitted,
		SubmittedAt:        &now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.Create(ctx, ref); err != nil {
		return nil, err
	}
	return ref, nil
}

func (s *achievementReferenceService) GetByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *achievementReferenceService) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error) {
	if studentID == uuid.Nil {
		return nil, errors.New("studentId is required")
	}
	return s.repo.GetByStudentID(ctx, studentID)
}

func (s *achievementReferenceService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.AchievementStatus, verifiedBy *uuid.UUID, rejectionNote *string) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}
	var verifiedAt *time.Time
	now := time.Now()

	switch status {
	case models.StatusVerified:
		verifiedAt = &now
	case models.StatusRejected:
		// when rejected, set verifiedAt to now as record of decision time
		verifiedAt = &now
	default:
		// for other statuses leave verifiedAt nil
	}

	return s.repo.UpdateStatus(ctx, id, status, verifiedBy, verifiedAt, rejectionNote)
}

func (s *achievementReferenceService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}
