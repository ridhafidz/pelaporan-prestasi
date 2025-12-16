package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"backend/app/models"
	"backend/app/repository"
)

// AchievementService defines business logic for achievements
type AchievementService interface {
	Create(ctx context.Context, a *models.Achievement) (*models.Achievement, error)
	GetByID(ctx context.Context, id string) (*models.Achievement, error)
	GetByStudentID(ctx context.Context, studentID string) ([]models.Achievement, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	SoftDelete(ctx context.Context, id string) error
	AddAttachment(ctx context.Context, id string, att models.Attachment) error
	RemoveAttachment(ctx context.Context, id string, fileName string) error
	FindByTag(ctx context.Context, tag string, limit int64) ([]models.Achievement, error)
}

type achievementService struct {
	repo repository.AchievementRepository
}

// NewAchievementService creates a new AchievementService
func NewAchievementService(repo repository.AchievementRepository) AchievementService {
	return &achievementService{repo: repo}
}

func (s *achievementService) Create(ctx context.Context, a *models.Achievement) (*models.Achievement, error) {
	if a == nil {
		return nil, errors.New("achievement is nil")
	}
	if a.StudentID == "" {
		return nil, errors.New("studentId is required")
	}
	if a.Title == "" {
		return nil, errors.New("title is required")
	}

	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *achievementService) GetByID(ctx context.Context, id string) (*models.Achievement, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *achievementService) GetByStudentID(ctx context.Context, studentID string) ([]models.Achievement, error) {
	if studentID == "" {
		return nil, errors.New("studentId is required")
	}
	return s.repo.GetByStudentID(ctx, studentID)
}

func (s *achievementService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if id == "" {
		return errors.New("id is required")
	}
	if updates == nil || len(updates) == 0 {
		return errors.New("updates is empty")
	}
	// prevent changing id or createdAt
	delete(updates, "_id")
	delete(updates, "id")
	delete(updates, "createdAt")

	// build $set
	set := bson.M{}
	for k, v := range updates {
		set[k] = v
	}
	// updatedAt
	set["updatedAt"] = time.Now()

	return s.repo.Update(ctx, id, bson.M{"$set": set})
}

func (s *achievementService) SoftDelete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *achievementService) AddAttachment(ctx context.Context, id string, att models.Attachment) error {
	if id == "" {
		return errors.New("id is required")
	}
	if att.FileName == "" || att.FileURL == "" {
		return errors.New("attachment fileName and fileUrl are required")
	}
	att.UploadedAt = time.Now()
	return s.repo.AddAttachment(ctx, id, att)
}

func (s *achievementService) RemoveAttachment(ctx context.Context, id string, fileName string) error {
	if id == "" || fileName == "" {
		return errors.New("id and fileName are required")
	}
	return s.repo.RemoveAttachmentByFileName(ctx, id, fileName)
}

func (s *achievementService) FindByTag(ctx context.Context, tag string, limit int64) ([]models.Achievement, error) {
	if tag == "" {
		return nil, errors.New("tag is required")
	}
	return s.repo.FindByTag(ctx, tag, limit)
}
