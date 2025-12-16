package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"fullName" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type UserData struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"fullName"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token        string   `json:"token"`
		RefreshToken string   `json:"refreshToken"`
		User         UserData `json:"user"`
	} `json:"data"`
}

type RegisterResponse struct {
	Status string `json:"status"`
	Data   struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"data"`
}

type RefreshTokenResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
}

type LogoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserClaims struct {
	UserID      uuid.UUID `json:"userId"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
