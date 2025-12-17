package service

import (
	"context"
	"errors"

	"backend/app/models"
	"backend/app/repository"

	"github.com/google/uuid"
)

type AchievementReferenceService interface {
	Create(ctx context.Context, studentID uuid.UUID, mongoAchievementID string) (*models.AchievementReference, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error)
	GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]models.AchievementReference, error)
	Submit(ctx context.Context, mongoID string) error
	Verify(ctx context.Context, mongoID string, verifierID uuid.UUID) error
	Reject(ctx context.Context, mongoID string, note string) error
	Delete(ctx context.Context, mongoID string) error
}

type achievementReferenceService struct {
	repo            repository.AchievementReferenceRepository
	achievementRepo repository.AchievementRepository
}

func NewAchievementReferenceService(
	repo repository.AchievementReferenceRepository,
	achievementRepo repository.AchievementRepository,
) AchievementReferenceService {
	return &achievementReferenceService{
		repo:            repo,
		achievementRepo: achievementRepo,
	}
}

func (s *achievementReferenceService) Create(
	ctx context.Context,
	studentID uuid.UUID,
	mongoAchievementID string,
) (*models.AchievementReference, error) {

	existing, err := s.repo.GetByMongoID(ctx, mongoAchievementID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("achievement reference already exists")
	}

	ref := &models.AchievementReference{
		StudentID:          studentID,
		MongoAchievementID: mongoAchievementID,
	}

	if err := s.repo.Create(ctx, ref); err != nil {
		return nil, err
	}

	return ref, nil
}

func (s *achievementReferenceService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.AchievementReference, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *achievementReferenceService) GetByMongoID(
	ctx context.Context,
	mongoID string,
) (*models.AchievementReference, error) {
	return s.repo.GetByMongoID(ctx, mongoID)
}

func (s *achievementReferenceService) GetByStudentID(
	ctx context.Context,
	studentID uuid.UUID,
	limit, offset int,
) ([]models.AchievementReference, error) {
	return s.repo.GetByStudentID(ctx, studentID, limit, offset)
}

func (s *achievementReferenceService) Submit(
	ctx context.Context,
	mongoID string,
) error {

	ref, err := s.repo.GetByMongoID(ctx, mongoID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}

	if ref.Status != models.StatusDraft {
		return errors.New("only draft achievement can be submitted")
	}

	return s.repo.UpdateStatus(ctx, mongoID, models.StatusSubmitted)
}

func (s *achievementReferenceService) Verify(
	ctx context.Context,
	mongoID string,
	verifierID uuid.UUID,
) error {

	ref, err := s.repo.GetByMongoID(ctx, mongoID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}

	if ref.Status != models.StatusSubmitted {
		return errors.New("only submitted achievement can be verified")
	}

	return s.repo.Verify(ctx, mongoID, verifierID)
}

func (s *achievementReferenceService) Reject(
	ctx context.Context,
	mongoID string,
	note string,
) error {

	if note == "" {
		return errors.New("rejection note is required")
	}

	ref, err := s.repo.GetByMongoID(ctx, mongoID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}

	if ref.Status != models.StatusSubmitted {
		return errors.New("only submitted achievement can be rejected")
	}

	return s.repo.Reject(ctx, mongoID, note)
}

func (s *achievementReferenceService) Delete(
	ctx context.Context,
	mongoID string,
) error {

	ref, err := s.repo.GetByMongoID(ctx, mongoID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement not found")
	}

	if ref.Status != models.StatusDraft {
		return errors.New("only draft achievement can be deleted")
	}

	if err := s.achievementRepo.SoftDelete(ctx, mongoID); err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, mongoID, models.StatusDeleted)
}
