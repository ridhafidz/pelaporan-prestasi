package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"fullName" db:"full_name"`
	RoleID       uuid.UUID `json:"roleId" db:"role_id"`
	RoleName     string    `json:"role" db:"role_name"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"fullName"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	IsActive    bool      `json:"isActive"`
}

type CreateUserRequest struct {
	Username string    `json:"username" validate:"required,max=50"`
	Email    string    `json:"email" validate:"required,email,max=100"`
	FullName string    `json:"fullName" validate:"required,max=100"`
	Password string    `json:"password" validate:"required,min=6"`
	RoleID   uuid.UUID `json:"roleId" validate:"required"`

	StudentID    string `json:"studentId,omitempty"`
	ProgramStudy string `json:"programStudy,omitempty"`
	AcademicYear string `json:"academicYear,omitempty"`

	LecturerID string `json:"lecturerId,omitempty"`
	Department string `json:"department,omitempty"`
}

type UpdateUserRequest struct {
	FullName *string    `json:"fullName" validate:"omitempty,max=100"`
	Username *string    `json:"username" validate:"omitempty,max=50"`
	Email    *string    `json:"email" validate:"omitempty,email,max=100"`
	RoleID   *uuid.UUID `json:"roleId"`
	IsActive *bool      `json:"isActive"`
}

type UpdateUserRoleRequest struct {
	RoleID uuid.UUID `json:"roleId" validate:"required"`
}
