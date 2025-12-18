package service

import (
	"context"
	"errors"
	"time"

	"backend/app/models"
	"backend/app/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAllUsers(ctx context.Context, page, limit int) ([]models.UserResponse, error)
	UpdateUserRole(ctx context.Context, id uuid.UUID, req *models.UpdateUserRoleRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(
	ctx context.Context,
	req *models.CreateUserRequest,
) (*models.UserResponse, error) {

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	userID := uuid.New()

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	newUser := &models.User{
		ID:           userID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.CreateTx(ctx, tx, newUser); err != nil {
		return nil, err
	}

	roleName, err := s.userRepo.GetRoleNameByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}

	switch roleName {

	case "Mahasiswa":
		if req.StudentID == "" || req.ProgramStudy == "" || req.AcademicYear == "" {
			return nil, errors.New("student data is required for role Mahasiswa")
		}
		err = s.userRepo.CreateStudentTx(ctx, tx, &models.Student{
			ID:           uuid.New(),
			UserID:       userID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})

	case "DosenWali":
		if req.LecturerID == "" || req.Department == "" {
			return nil, errors.New("lecturer data is required for role DosenWali")
		}
		err = s.userRepo.CreateLecturerTx(ctx, tx, &models.Lecturer{
			ID:         uuid.New(),
			UserID:     userID,
			LecturerID: req.LecturerID,
			Department: req.Department,
			CreatedAt:  time.Now(),
		})
	}

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	createdUser, _ := s.userRepo.FindByID(ctx, userID)
	return s.mapToResponse(ctx, createdUser)
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) mapToResponse(
	ctx context.Context,
	user *models.User,
) (*models.UserResponse, error) {

	perms, err := s.userRepo.GetPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        user.RoleName,
		IsActive:    user.IsActive,
		Permissions: perms,
	}, nil
}

func (s *userService) GetAllUsers(ctx context.Context, page, limit int) ([]models.UserResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, err := s.userRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var userResponses []models.UserResponse
	for _, u := range users {
		userTemp := u
		resp, err := s.mapToResponse(ctx, &userTemp)
		if err != nil {
			return nil, err
		}
		userResponses = append(userResponses, *resp)
	}

	return userResponses, nil
}

func (s *userService) UpdateUserRole(ctx context.Context, id uuid.UUID, req *models.UpdateUserRoleRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	user.RoleID = req.RoleID
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}
