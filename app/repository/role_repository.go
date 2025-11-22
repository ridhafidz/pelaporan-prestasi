package repository

import (
	"context"
	"database/sql"

	"backend/app/models"
)

type IRoleRepository interface {
	GetAll(ctx context.Context) ([]models.Role, error)
}

type RoleRepository struct {
	DB *sql.DB
}

func NewRoleRepository(db *sql.DB) IRoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]models.Role, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, name, description, created_at FROM roles
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		roles = append(roles, role)
	}
	return roles, nil
}
