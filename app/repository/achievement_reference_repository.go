package repository

import (
	"context"
	"database/sql"

	"backend/app/models"
)

type IAchievementReferenceRepo interface {
	Create(ctx context.Context, ar *models.AchievementReference) error
	UpdateStatus(ctx context.Context, id string, status string) error
	FindByID(ctx context.Context, id string) (*models.AchievementReference, error)
	GetByStudent(ctx context.Context, studentID string) ([]models.AchievementReference, error)
}

type AchievementReferenceRepository struct {
	DB *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) IAchievementReferenceRepo {
	return &AchievementReferenceRepository{DB: db}
}

func (r *AchievementReferenceRepository) Create(ctx context.Context, ar *models.AchievementReference) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, ar.ID, ar.StudentID, ar.MongoAchievementID, ar.Status)
	return err
}

func (r *AchievementReferenceRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE achievement_references SET status=$1, updated_at=NOW() WHERE id=$2
	`, status, id)
	return err
}

func (r *AchievementReferenceRepository) FindByID(ctx context.Context, id string) (*models.AchievementReference, error) {
	var a models.AchievementReference
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references WHERE id = $1
	`, id).Scan(
		&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.SubmittedAt,
		&a.VerifiedAt, &a.VerifiedBy, &a.RejectionNote, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AchievementReferenceRepository) GetByStudent(ctx context.Context, studentID string) ([]models.AchievementReference, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references WHERE student_id=$1
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementReference
	for rows.Next() {
		var a models.AchievementReference
		rows.Scan(
			&a.ID, &a.StudentID, &a.MongoAchievementID, &a.Status, &a.SubmittedAt,
			&a.VerifiedAt, &a.VerifiedBy, &a.RejectionNote, &a.CreatedAt, &a.UpdatedAt,
		)
		list = append(list, a)
	}
	return list, nil
}
