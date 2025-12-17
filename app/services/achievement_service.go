package service

import (
	"context"
	"errors"
	"time"

	"backend/app/models"
	"backend/app/repository"
)

type AchievementService interface {
	Create(ctx context.Context, achievement *models.Achievement) (string, error)
	GetByID(ctx context.Context, id string) (*models.Achievement, error)
	AddAttachment(ctx context.Context, id string, attachment models.Attachment) error
}

type achievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(
	repo repository.AchievementRepository,
) AchievementService {
	return &achievementService{
		repo: repo,
	}
}

func (s *achievementService) Create(
	ctx context.Context,
	achievement *models.Achievement,
) (string, error) {

	if achievement == nil {
		return "", errors.New("achievement payload is required")
	}

	// Validasi minimal (bisa kamu perluas)
	if achievement.Title == "" {
		return "", errors.New("achievement title is required")
	}

	if achievement.StudentID == "" {
		return "", errors.New("student_id is required")
	}

	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	return s.repo.Create(ctx, achievement)
}

func (s *achievementService) GetByID(
	ctx context.Context,
	id string,
) (*models.Achievement, error) {

	if id == "" {
		return nil, errors.New("achievement id is required")
	}

	return s.repo.FindByID(ctx, id)
}

func (s *achievementService) AddAttachment(
	ctx context.Context,
	id string,
	attachment models.Attachment,
) error {

	if id == "" {
		return errors.New("achievement id is required")
	}

	if attachment.FileURL == "" {
		return errors.New("attachment file_url is required")
	}

	attachment.UploadedAt = time.Now()

	return s.repo.AddAttachment(ctx, id, attachment)
}
