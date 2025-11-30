package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"backend/app/models"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, userID string, roleID string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password_hash, full_name, role_id,
		       is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // NOT FOUND
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name,
		                   role_id, is_active, created_at, updated_at)
		VALUES (:id, :username, :email, :password_hash, :full_name,
		        :role_id, :is_active, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = :username,
		    email = :email,
		    full_name = :full_name,
		    role_id = :role_id,
		    is_active = :is_active,
		    updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) UpdateRole(ctx context.Context, userID string, roleID string) error {
	query := `UPDATE users SET role_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, roleID, userID)
	return err
}

