package repository

import (
	"context"
	"database/sql"

	"backend/app/models"
)

type IUserRepository interface {
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]models.User, error)
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive,
	)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users SET username=$1, email=$2, full_name=$3, role_id=$4, is_active=$5, updated_at=NOW()
		WHERE id=$6
	`,
		user.Username, user.Email, user.FullName, user.RoleID, user.IsActive, user.ID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash,
			&u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		)
		list = append(list, u)
	}
	return list, nil
}
