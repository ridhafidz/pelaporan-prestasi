package models

import (
	"time"

	"github.com/google/uuid"
)

type Lecturer struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"userId" db:"user_id"`
	LecturerID string    `json:"lecturerId" db:"lecturer_id"`
	Department string    `json:"department" db:"department"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

type CreateLecturerRequest struct {
	UserID     uuid.UUID `json:"userId" validate:"required"`
	LecturerID string    `json:"lecturerId" validate:"required,max=20"`
	Department string    `json:"department" validate:"required,max=100"`
}

type LecturerDetailResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	FullName   string    `json:"fullName"`   // Diambil dari tabel users
	LecturerID string    `json:"lecturerId"`
	Department string    `json:"department"`
}