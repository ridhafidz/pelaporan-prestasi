package service

import (
    "context"
    "errors"
    "time"
    "backend/app/models"     
    "backend/app/repository" 

    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
    RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.RefreshTokenResponse, error)
    Logout(ctx context.Context, req models.LogoutRequest) error
    GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserData, error)
}

type authService struct {
    repo      repository.AuthRepository
    jwtSecret []byte
}

func NewAuthService(repo repository.AuthRepository, secret string) AuthService {
    return &authService{
        repo:      repo,
        jwtSecret: []byte(secret),
    }
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
    user, err := s.repo.FindByUsername(ctx, req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
 
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
 
    permissions, err := s.repo.GetPermissionsByRoleID(ctx, user.RoleID)
    if err != nil {
        return nil, err
    }

    accessToken, err := s.generateToken(user, permissions, 15*time.Minute) // Short lived
    if err != nil {
        return nil, err
    }

    refreshTokenStr, err := s.generateToken(user, permissions, 7*24*time.Hour) // Long lived
    if err != nil {
        return nil, err
    }

    refTokenModel := models.RefreshToken{
        ID:        uuid.New(),
        UserID:    user.ID,
        Token:     refreshTokenStr,
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
        CreatedAt: time.Now(),
    }
    if err := s.repo.StoreRefreshToken(ctx, refTokenModel); err != nil {
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
            RefreshToken: refreshTokenStr,
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

func (s *authService) RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
    storedToken, err := s.repo.GetRefreshToken(ctx, req.RefreshToken)
    if err != nil {
        return nil, errors.New("invalid refresh token")
    }

    if time.Now().After(storedToken.ExpiresAt) {
        return nil, errors.New("refresh token expired")
    }
    user, err := s.repo.FindByID(ctx, storedToken.UserID)
    if err != nil {
        return nil, err
    }
    permissions, _ := s.repo.GetPermissionsByRoleID(ctx, user.RoleID)

    newAccessToken, err := s.generateToken(user, permissions, 15*time.Minute)
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
            RefreshToken: req.RefreshToken, 
        },
    }, nil
}

func (s *authService) Logout(ctx context.Context, req models.LogoutRequest) error {
    return s.repo.DeleteRefreshToken(ctx, req.RefreshToken)
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserData, error) {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    permissions, _ := s.repo.GetPermissionsByRoleID(ctx, user.RoleID)

    return &models.UserData{
        ID:          user.ID,
        Username:    user.Username,
        FullName:    user.FullName,
        Role:        user.RoleName,
        Permissions: permissions,
    }, nil
}

func (s *authService) generateToken(user *models.User, permissions []string, ttl time.Duration) (string, error) {
    claims := models.UserClaims{
        UserID:      user.ID,
        Username:    user.Username,
        Role:        user.RoleName,
        Permissions: permissions,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "sistem-prestasi-mahasiswa",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}