package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"` 
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// CreateRoleRequest digunakan untuk validasi input saat membuat role baru
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}

// UpdateRoleRequest digunakan untuk validasi input saat update role
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}