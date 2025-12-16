package repository

import (
	"context"
	"database/sql"

	"backend/app/models"

	"github.com/jmoiron/sqlx"
)

type StudentRepository interface {
	GetAll(ctx context.Context) ([]models.Student, error)
	GetByID(ctx context.Context, id string) (*models.Student, error)
	Create(ctx context.Context, student *models.Student) error
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id string) error
}

type studentRepository struct {
	db *sqlx.DB
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) GetAll(ctx context.Context) ([]models.Student, error) {
	var students []models.Student

	query := `
        SELECT id, name, email, enrollment_year, is_active, created_at, updated_at
        FROM students
        ORDER BY created_at DESC
    `

	err := r.db.SelectContext(ctx, &students, query)
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *studentRepository) GetByID(ctx context.Context, id string) (*models.Student, error) {
	var student models.Student

	query := `
        SELECT id, name, email, enrollment_year, is_active, created_at, updated_at
        FROM students
        WHERE id = $1
    `

	err := r.db.GetContext(ctx, &student, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &student, nil
}

func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
	query := `
        INSERT INTO students (name, email, enrollment_year, is_active, created_at, updated_at)
        VALUES (:name, :email, :enrollment_year, :is_active, :created_at, :updated_at)
        RETURNING id
    `

	// Use NamedQuery to capture returning id
	rows, err := r.db.NamedQueryContext(ctx, query, student)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&student.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
	query := `
        UPDATE students
        SET name = :name,
            email = :email,
            enrollment_year = :enrollment_year,
            is_active = :is_active,
            updated_at = :updated_at
        WHERE id = :id
    `

	_, err := r.db.NamedExecContext(ctx, query, student)
	return err
}

func (r *studentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM students WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
