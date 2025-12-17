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
	Update(ctx context.Context, id string, achievement *models.Achievement) error
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
func (s *achievementService) Update(
	ctx context.Context,
	id string,
	payload *models.Achievement,
) error {

	if id == "" {
		return errors.New("achievement id is required")
	}
	if payload == nil {
		return errors.New("update payload is required")
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if payload.Title != "" {
		existing.Title = payload.Title
	}
	if payload.Description != "" {
		existing.Description = payload.Description
	}

	existing.UpdatedAt = time.Now()

	return s.repo.Update(ctx, id, existing)
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
