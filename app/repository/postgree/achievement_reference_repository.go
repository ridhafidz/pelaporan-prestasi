package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"backend/app/models"
)

// AchievementReferenceRepository defines operations for achievement references stored in PostgreSQL
type AchievementReferenceRepository interface {
	GetAll(ctx context.Context) ([]models.AchievementReference, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error)
	Create(ctx context.Context, ref *models.AchievementReference) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.AchievementStatus, verifiedBy *uuid.UUID, verifiedAt *time.Time, rejectionNote *string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type achievementReferenceRepo struct {
	db *sqlx.DB
}

func NewAchievementReferenceRepository(db *sqlx.DB) AchievementReferenceRepository {
	return &achievementReferenceRepo{db: db}
}

func (r *achievementReferenceRepo) GetAll(ctx context.Context) ([]models.AchievementReference, error) {
	var refs []models.AchievementReference
	query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        ORDER BY created_at DESC
    `
	if err := r.db.SelectContext(ctx, &refs, query); err != nil {
		return nil, err
	}
	return refs, nil
}

func (r *achievementReferenceRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error) {
	var ref models.AchievementReference
	query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE id = $1
    `
	err := r.db.GetContext(ctx, &ref, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ref, nil
}

func (r *achievementReferenceRepo) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.AchievementReference, error) {
	var refs []models.AchievementReference
	query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE student_id = $1
        ORDER BY created_at DESC
    `
	if err := r.db.SelectContext(ctx, &refs, query, studentID); err != nil {
		return nil, err
	}
	return refs, nil
}

func (r *achievementReferenceRepo) Create(ctx context.Context, ref *models.AchievementReference) error {
	if ref == nil {
		return sql.ErrNoRows
	}
	now := time.Now()
	if ref.ID == uuid.Nil {
		ref.ID = uuid.New()
	}
	if ref.CreatedAt.IsZero() {
		ref.CreatedAt = now
	}
	ref.UpdatedAt = now

	query := `
        INSERT INTO achievement_references (id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
        VALUES (:id, :student_id, :mongo_achievement_id, :status, :submitted_at, :verified_at, :verified_by, :rejection_note, :created_at, :updated_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, ref)
	return err
}

func (r *achievementReferenceRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.AchievementStatus, verifiedBy *uuid.UUID, verifiedAt *time.Time, rejectionNote *string) error {
	query := `
        UPDATE achievement_references
        SET status = $1,
            verified_by = $2,
            verified_at = $3,
            rejection_note = $4,
            updated_at = $5
        WHERE id = $6
    `
	res, err := r.db.ExecContext(ctx, query, status, verifiedBy, verifiedAt, rejectionNote, time.Now(), id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *achievementReferenceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM achievement_references WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
