package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/app/models"

	"github.com/google/uuid"
)

type AchievementReferenceRepository interface {
	Create(ctx context.Context, ref *models.AchievementReference) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error)
	GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error) // Added this to interface
	GetByStudentID(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]models.AchievementReference, error)
	UpdateStatus(ctx context.Context, mongoID string, status models.AchievementStatus) error
	Verify(ctx context.Context, mongoID string, verifierID uuid.UUID) error
	Reject(ctx context.Context, mongoID string, note string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type achievementReferenceRepository struct {
	db *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) AchievementReferenceRepository {
	return &achievementReferenceRepository{db: db}
}

func (r *achievementReferenceRepository) Create(ctx context.Context, ref *models.AchievementReference) error {
	query := `
        INSERT INTO achievement_references (
            id, student_id, mongo_achievement_id, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `
	ref.ID = uuid.New()
	ref.Status = models.StatusDraft
	now := time.Now()
	ref.CreatedAt = now
	ref.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

func (r *achievementReferenceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.AchievementReference, error) {
	query := `
        SELECT id, student_id, mongo_achievement_id, status, 
               submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE id = $1
    `
	return r.scanRow(ctx, query, id)
}

func (r *achievementReferenceRepository) GetByMongoID(ctx context.Context, mongoID string) (*models.AchievementReference, error) {
	query := `
        SELECT id, student_id, mongo_achievement_id, status, 
               submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE mongo_achievement_id = $1
    `
	return r.scanRow(ctx, query, mongoID)
}

func (r *achievementReferenceRepository) scanRow(ctx context.Context, query string, arg interface{}) (*models.AchievementReference, error) {
	var ref models.AchievementReference

	// Variables to handle SQL NULLs
	var submittedAt, verifiedAt sql.NullTime
	var verifiedBy sql.NullString
	var rejectionNote sql.NullString

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
		&submittedAt, &verifiedAt, &verifiedBy, &rejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}

	// Map Nullable types back to struct
	if submittedAt.Valid {
		ref.SubmittedAt = &submittedAt.Time
	}
	if verifiedAt.Valid {
		ref.VerifiedAt = &verifiedAt.Time
	}
	if verifiedBy.Valid && verifiedBy.String != "" {
		uid := uuid.MustParse(verifiedBy.String)
		ref.VerifiedBy = &uid
	}
	if rejectionNote.Valid {
		ref.RejectionNote = &rejectionNote.String
	}

	return &ref, nil
}

func (r *achievementReferenceRepository) UpdateStatus(ctx context.Context, mongoID string, status models.AchievementStatus) error {
	query := `
        UPDATE achievement_references 
        SET status = $1, updated_at = NOW(), submitted_at = CASE WHEN $1 = 'submitted' THEN NOW() ELSE submitted_at END
        WHERE mongo_achievement_id = $2
    `
	_, err := r.db.ExecContext(ctx, query, status, mongoID)
	return err
}

func (r *achievementReferenceRepository) Verify(ctx context.Context, mongoID string, verifierID uuid.UUID) error {
	query := `
        UPDATE achievement_references 
        SET status = 'verified', verified_at = NOW(), verified_by = $1, updated_at = NOW()
        WHERE mongo_achievement_id = $2
    `
	_, err := r.db.ExecContext(ctx, query, verifierID, mongoID)
	return err
}

func (r *achievementReferenceRepository) Reject(ctx context.Context, mongoID string, note string) error {
	query := `
        UPDATE achievement_references 
        SET status = 'rejected', rejection_note = $1, updated_at = NOW()
        WHERE mongo_achievement_id = $2
    `
	_, err := r.db.ExecContext(ctx, query, note, mongoID)
	return err
}

func (r *achievementReferenceRepository) GetByStudentID(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]models.AchievementReference, error) {
	query := `
        SELECT id, mongo_achievement_id, status, updated_at
        FROM achievement_references
        WHERE student_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		if err := rows.Scan(&ref.ID, &ref.MongoAchievementID, &ref.Status, &ref.UpdatedAt); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *achievementReferenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM achievement_references WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
