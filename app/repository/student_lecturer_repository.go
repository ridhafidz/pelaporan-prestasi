package repository

import (
	"context"
	"database/sql"

	"backend/app/models"

	"github.com/google/uuid"
)

type StudentLecturerRepository interface {
	GetAllStudents(ctx context.Context) ([]models.StudentDetailResponse, error)
	GetStudentByID(ctx context.Context, id uuid.UUID) (*models.StudentDetailResponse, error)
	GetStudentAchievements(ctx context.Context, studentID uuid.UUID) ([]models.StudentAchievementResponse, error)
	UpdateStudentAdvisor(ctx context.Context, studentID uuid.UUID, advisorID *uuid.UUID) error

	GetAllLecturers(ctx context.Context) ([]models.LecturerDetailResponse, error)
	GetLecturerAdvisees(ctx context.Context, lecturerID uuid.UUID) ([]models.StudentDetailResponse, error)
	GetLecturerByUserID(ctx context.Context, userID uuid.UUID) (*models.LecturerDetailResponse, error)
	IsAdvisorOf(ctx context.Context, advisorUserID uuid.UUID, studentID uuid.UUID) (bool, error)
}

type studentLecturerRepository struct {
	db *sql.DB
}

func NewStudentLecturerRepository(db *sql.DB) StudentLecturerRepository {
	return &studentLecturerRepository{db: db}
}

func (r *studentLecturerRepository) GetAllStudents(
	ctx context.Context,
) ([]models.StudentDetailResponse, error) {

	query := `
		SELECT 
			s.id,
			s.user_id,
			u.full_name,
			s.student_id,
			s.program_study,
			s.academic_year,
			COALESCE(adv.full_name, '') AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users adv ON adv.id = l.user_id
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.StudentDetailResponse
	for rows.Next() {
		var s models.StudentDetailResponse
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.FullName,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorName,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}
func (r *studentLecturerRepository) GetStudentByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.StudentDetailResponse, error) {

	query := `
		SELECT 
			s.id,
			s.user_id,
			u.full_name,
			s.student_id,
			s.program_study,
			s.academic_year,
			COALESCE(adv.full_name, '') AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users adv ON adv.id = l.user_id
		WHERE s.id = $1
	`

	var s models.StudentDetailResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.FullName,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *studentLecturerRepository) GetStudentAchievements(
	ctx context.Context,
	studentID uuid.UUID,
) ([]models.StudentAchievementResponse, error) {

	query := `
		SELECT
			sar.mongo_achievement_id,
			sar.title,
			sar.status,
			sar.submitted_at,
			sar.verified_at
		FROM student_achievement_references sar
		WHERE sar.student_id = $1
		ORDER BY sar.submitted_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.StudentAchievementResponse

	for rows.Next() {
		var a models.StudentAchievementResponse
		if err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Status,
			&a.SubmittedAt,
			&a.VerifiedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

func (r *studentLecturerRepository) UpdateStudentAdvisor(
	ctx context.Context,
	studentID uuid.UUID,
	advisorID *uuid.UUID,
) error {

	query := `
		UPDATE students
		SET advisor_id = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, advisorID, studentID)
	return err
}
func (r *studentLecturerRepository) GetAllLecturers(
	ctx context.Context,
) ([]models.LecturerDetailResponse, error) {

	query := `
		SELECT
			l.id,
			l.user_id,
			u.full_name,
			l.lecturer_id,
			l.department
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY u.full_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.LecturerDetailResponse
	for rows.Next() {
		var l models.LecturerDetailResponse
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.FullName,
			&l.LecturerID,
			&l.Department,
		); err != nil {
			return nil, err
		}
		result = append(result, l)
	}

	return result, nil
}
func (r *studentLecturerRepository) GetLecturerAdvisees(
	ctx context.Context,
	lecturerID uuid.UUID,
) ([]models.StudentDetailResponse, error) {

	query := `
		SELECT 
			s.id,
			s.user_id,
			u.full_name,
			s.student_id,
			s.program_study,
			s.academic_year,
			u2.full_name AS advisor_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		JOIN lecturers l ON l.id = s.advisor_id
		JOIN users u2 ON u2.id = l.user_id
		WHERE l.id = $1
		ORDER BY u.full_name
	`

	rows, err := r.db.QueryContext(ctx, query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.StudentDetailResponse
	for rows.Next() {
		var s models.StudentDetailResponse
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.FullName,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorName,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (r *studentLecturerRepository) GetLecturerByUserID(ctx context.Context, userID uuid.UUID) (*models.LecturerDetailResponse, error) {
	query := `
        SELECT l.id, l.user_id, u.full_name, l.lecturer_id, l.department
        FROM lecturers l
        JOIN users u ON u.id = l.user_id
        WHERE l.user_id = $1`

	var l models.LecturerDetailResponse
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&l.ID, &l.UserID, &l.FullName, &l.LecturerID, &l.Department)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *studentLecturerRepository) IsAdvisorOf(ctx context.Context, advisorUserID uuid.UUID, studentID uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 FROM students s
            JOIN lecturers l ON s.advisor_id = l.id
            WHERE l.user_id = $1 AND s.id = $2
        )`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, advisorUserID, studentID).Scan(&exists)
	return exists, err
}
