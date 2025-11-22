package repository

import (
	"context"
	"database/sql"

	"backend/app/models"
)

type IStudentRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.Student, error)
	GetAdvisees(ctx context.Context, advisorID string) ([]models.Student, error)
}

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) IStudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) FindByUserID(ctx context.Context, userID string) (*models.Student, error) {
	var s models.Student
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students WHERE user_id = $1
	`, userID).Scan(
		&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear,
		&s.AdvisorID, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) GetAdvisees(ctx context.Context, advisorID string) ([]models.Student, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students WHERE advisor_id = $1
	`, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student
	for rows.Next() {
		var s models.Student
		rows.Scan(
			&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear,
			&s.AdvisorID, &s.CreatedAt,
		)
		list = append(list, s)
	}
	return list, nil
}
