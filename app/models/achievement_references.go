package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
)

type AchievementReference struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	StudentID          uuid.UUID         `json:"student_id" db:"student_id"`
	MongoAchievementID string            `json:"mongo_achievement_id" db:"mongo_achievement_id"`
	Status             AchievementStatus `json:"status" db:"status"`

	SubmittedAt  sql.NullTime   `json:"submittedAt,omitempty" db:"submitted_at"`
	VerifiedAt   sql.NullTime   `json:"verifiedAt,omitempty" db:"verified_at"`
	VerifiedBy   *uuid.UUID     `json:"verifiedBy,omitempty" db:"verified_by"`
	RejectionNote sql.NullString `json:"rejectionNote,omitempty" db:"rejection_note"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
