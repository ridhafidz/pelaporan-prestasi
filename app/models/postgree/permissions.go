package models

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"` // UUID
	Name        string    `json:"name" db:"name"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description" db:"description"`
}

type PermissionResponse struct {
	ID          uuid.UUID `json:"id" db:"id"` // UUID
	Name        string    `json:"name" db:"name"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description" db:"description"`
	CreatedAt   string    `json:"created_at" db:"created_at"`
	UpdatedAt   string    `json:"updated_at" db:"updated_at"`
}
