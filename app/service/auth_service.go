package service

import (
	"context"
	"errors"
	"time"

	"backend/app/models"
	"backend/app/repository"
	"backend/app/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error)
	Logout(ctx context.Context, req models.LogoutRequest) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserData, error)
}

type authService struct {
	authrepo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{authrepo: repo}
}

func (s *authService) Login(
	ctx context.Context,
	req models.LoginRequest,
) (*models.LoginResponse, error) {

	user, err := s.authrepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	permissions, err := s.authrepo.GetPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	claims := &models.JWTClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Role:        user.RoleName,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistem-prestasi-mahasiswa",
		},
	}

	refreshTokenClaims := &models.JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistem-prestasi-mahasiswa",
		},
	}

	refreshToken, err := utils.GenerateTokenWithClaims(refreshTokenClaims)
	if err != nil {
		return nil, err
	}

	switch user.RoleName {
	case "Mahasiswa":
		studentID, err := s.authrepo.GetStudentIDByUserID(ctx, user.ID)
		if err != nil {
			return nil, errors.New("student profile not found")
		}
		claims.StudentID = &studentID

	case "DosenWali":
		lecturerID, err := s.authrepo.GetLecturerIDByUserID(ctx, user.ID)
		if err != nil {
			return nil, errors.New("lecturer profile not found")
		}
		claims.LecturerID = &lecturerID
	}

	accessToken, err := utils.GenerateTokenWithClaims(claims)
	if err != nil {
		return nil, err
	}

	refToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.authrepo.StoreRefreshToken(ctx, refToken); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Status: "success",
		Data: struct {
			Token        string          `json:"token"`
			RefreshToken string          `json:"refreshToken"`
			User         models.UserData `json:"user"`
		}{
			Token:        accessToken,
			RefreshToken: refreshToken,
			User: models.UserData{
				ID:          user.ID,
				Username:    user.Username,
				FullName:    user.FullName,
				Role:        user.RoleName,
				Permissions: permissions,
			},
		},
	}, nil
}

func (s *authService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*models.RefreshTokenResponse, error) {

	stored, err := s.authrepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.authrepo.FindByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}

	permissions, err := s.authrepo.GetPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	claims := &models.JWTClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Role:        user.RoleName,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sistem-prestasi-mahasiswa",
		},
	}

	switch user.RoleName {
	case "Mahasiswa":
		studentID, err := s.authrepo.GetStudentIDByUserID(ctx, user.ID)
		if err != nil {
			return nil, errors.New("student profile not found")
		}
		claims.StudentID = &studentID

	case "DosenWali":
		lecturerID, err := s.authrepo.GetLecturerIDByUserID(ctx, user.ID)
		if err != nil {
			return nil, errors.New("lecturer profile not found")
		}
		claims.LecturerID = &lecturerID
	}

	newAccessToken, err := utils.GenerateTokenWithClaims(claims)
	if err != nil {
		return nil, err
	}

	return &models.RefreshTokenResponse{
		Status: "success",
		Data: struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refreshToken"`
		}{
			Token:        newAccessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authService) Logout(
	ctx context.Context,
	req models.LogoutRequest,
) error {
	return s.authrepo.DeleteRefreshToken(ctx, req.RefreshToken)
}

func (s *authService) GetProfile(
	ctx context.Context,
	userID uuid.UUID,
) (*models.UserData, error) {

	user, err := s.authrepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissions, _ := s.authrepo.GetPermissionsByRoleID(ctx, user.RoleID)

	return &models.UserData{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        user.RoleName,
		Permissions: permissions,
	}, nil
}
