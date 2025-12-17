package service

import (
	"context"
	"errors"

	"backend/app/models"
	"backend/app/repository"

	"github.com/google/uuid"
)

type StudentLecturerService interface {
	GetStudents(ctx context.Context) ([]models.StudentDetailResponse, error)
	GetStudentDetail(ctx context.Context, studentID uuid.UUID) (*models.StudentDetailResponse, error)
	GetStudentAchievements(ctx context.Context, studentID uuid.UUID) ([]models.StudentAchievementResponse, error)
	UpdateAdvisor(ctx context.Context, studentID uuid.UUID, advisorID *uuid.UUID) error

	GetLecturers(ctx context.Context) ([]models.LecturerDetailResponse, error)
	GetLecturerAdvisees(ctx context.Context, lecturerID uuid.UUID) ([]models.StudentDetailResponse, error)
}
type studentLecturerService struct {
	repo repository.StudentLecturerRepository
}

func NewStudentLecturerService(
	repo repository.StudentLecturerRepository,
) StudentLecturerService {
	return &studentLecturerService{
		repo: repo,
	}
}
func (s *studentLecturerService) GetStudents(
	ctx context.Context,
) ([]models.StudentDetailResponse, error) {

	return s.repo.GetAllStudents(ctx)
}
func (s *studentLecturerService) GetStudentDetail(
	ctx context.Context,
	studentID uuid.UUID,
) (*models.StudentDetailResponse, error) {

	student, err := s.repo.GetStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, errors.New("student not found")
	}

	return student, nil
}

func (s *studentLecturerService) GetStudentAchievements(
	ctx context.Context,
	studentID uuid.UUID,
) ([]models.StudentAchievementResponse, error) {

	exists, err := s.repo.GetStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if exists == nil {
		return nil, errors.New("student not found")
	}

	return s.repo.GetStudentAchievements(ctx, studentID)
}

func (s *studentLecturerService) UpdateAdvisor(
	ctx context.Context,
	studentID uuid.UUID,
	advisorID *uuid.UUID,
) error {

	student, err := s.repo.GetStudentByID(ctx, studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	return s.repo.UpdateStudentAdvisor(ctx, studentID, advisorID)
}
func (s *studentLecturerService) GetLecturers(
	ctx context.Context,
) ([]models.LecturerDetailResponse, error) {

	return s.repo.GetAllLecturers(ctx)
}
func (s *studentLecturerService) GetLecturerAdvisees(
	ctx context.Context,
	lecturerID uuid.UUID,
) ([]models.StudentDetailResponse, error) {

	advisees, err := s.repo.GetLecturerAdvisees(ctx, lecturerID)
	if err != nil {
		return nil, err
	}

	if len(advisees) == 0 {
		return []models.StudentDetailResponse{}, nil
	}

	return advisees, nil
}
