package service

import (
	"context"
	"errors"
	"time"

	models "backend/app/models/postgree"
	repo "backend/app/repository/postgree"

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
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	existingUser, _ := s.userRepo.FindByUsernameOrEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, _ = s.userRepo.FindByUsernameOrEmail(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true, // Default active
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.FindByID(ctx, newUser.ID)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(createdUser), nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) error {
	// 1. Cari user lama
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

func (s *userService) mapToResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        user.RoleName,
		IsActive:    user.IsActive,
		Permissions: []string{},
	}
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
		userResponses = append(userResponses, *s.mapToResponse(&userTemp))
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
