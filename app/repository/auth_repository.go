package repository

import (
	"backend/app/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type AuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error)
	GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetLecturerIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	StoreRefreshToken(ctx context.Context, token models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
        SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, u.is_active, r.name as role_name
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.username = $1 AND u.is_active = true
    `
	var user models.User
	// Scan harus sesuai urutan query SELECT
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &user.RoleName,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	query := `
        SELECT p.name 
        FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = $1
    `
	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (r *authRepository) GetStudentIDByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (uuid.UUID, error) {

	query := `SELECT id FROM students WHERE user_id = $1`
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	return id, err
}

func (r *authRepository) GetLecturerIDByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (uuid.UUID, error) {

	query := `SELECT id FROM lecturers WHERE user_id = $1`
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	return id, err
}

func (r *authRepository) StoreRefreshToken(ctx context.Context, token models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *authRepository) GetRefreshToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1`
	var token models.RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *authRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT u.id, u.username, u.full_name, r.name as role_name, u.role_id
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.id = $1
    `
	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.FullName, &user.RoleName, &user.RoleID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
