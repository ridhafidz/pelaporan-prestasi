package models

import "github.com/google/uuid"

type RolePermission struct {
	RoleID       uuid.UUID `json:"roleId" db:"role_id"`            
	PermissionID uuid.UUID `json:"permissionId" db:"permission_id"`
}

// Struktur ini berguna jika nanti ada endpoint untuk menetapkan banyak permission ke satu role sekaligus
type AssignPermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permissionIds" validate:"required"`
}