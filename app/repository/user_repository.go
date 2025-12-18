package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend/app/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindAll(ctx context.Context, limit int, offset int) ([]models.User, error)
	FindByUsernameOrEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTx(ctx context.Context, tx *sql.Tx, user *models.User) error
	GetRoleNameByID(ctx context.Context, roleID uuid.UUID) (string, error)
	CreateStudentTx(ctx context.Context, tx *sql.Tx, student *models.Student) error
	CreateLecturerTx(ctx context.Context, tx *sql.Tx, lecturer *models.Lecturer) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (
            id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, 
            u.role_id, r.name as role_name, 
            u.is_active, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN roles r ON u.role_id = r.id
        WHERE u.id = $1
    `

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.RoleName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) FindByUsernameOrEmail(ctx context.Context, identity string) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, 
            u.role_id, r.name as role_name, 
            u.is_active, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN roles r ON u.role_id = r.id
        WHERE u.username = $1 OR u.email = $1
    `

	var user models.User
	err := r.db.QueryRowContext(ctx, query, identity).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.RoleName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetPermissionsByRoleID(
	ctx context.Context,
	roleID uuid.UUID,
) ([]string, error) {

	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
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

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET full_name = $1, username = $2, email = $3, role_id = $4, is_active = $5, updated_at = $6
        WHERE id = $7
    `
	_, err := r.db.ExecContext(ctx, query,
		user.FullName,
		user.Username,
		user.Email,
		user.RoleID,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) FindAll(ctx context.Context, limit int, offset int) ([]models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.full_name, 
            u.role_id, r.name as role_name, 
            u.is_active, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN roles r ON u.role_id = r.id
        LIMIT $1 OFFSET $2
    `
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.RoleID,
			&user.RoleName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *userRepository) CreateTx(
	ctx context.Context,
	tx *sql.Tx,
	user *models.User,
) error {

	query := `
		INSERT INTO users (
			id, username, email, password_hash,
			full_name, role_id, is_active,
			created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := tx.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetRoleNameByID(
	ctx context.Context,
	roleID uuid.UUID,
) (string, error) {

	query := `SELECT name FROM roles WHERE id = $1`

	var roleName string
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&roleName)
	return roleName, err
}

func (r *userRepository) CreateStudentTx(
	ctx context.Context,
	tx *sql.Tx,
	student *models.Student,
) error {

	query := `
		INSERT INTO students (
			id, user_id, student_id,
			program_study, academic_year,
			created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := tx.ExecContext(ctx, query,
		student.ID,
		student.UserID,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.CreatedAt,
		student.UpdatedAt,
	)

	return err
}

func (r *userRepository) CreateLecturerTx(
	ctx context.Context,
	tx *sql.Tx,
	lecturer *models.Lecturer,
) error {

	query := `
		INSERT INTO lecturers (
			id, user_id, lecturer_id,
			department, created_at
		)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := tx.ExecContext(ctx, query,
		lecturer.ID,
		lecturer.UserID,
		lecturer.LecturerID,
		lecturer.Department,
		lecturer.CreatedAt,
	)

	return err
}

