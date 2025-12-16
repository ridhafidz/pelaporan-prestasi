package models

import (
	"time"

	"github.com/google/uuid"
)

type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
	StatusDeleted   AchievementStatus = "deleted"
)

type AchievementReference struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	StudentID          uuid.UUID         `json:"studentId" db:"student_id"`
	MongoAchievementID string            `json:"mongoAchievementId" db:"mongo_achievement_id"`
	Status             AchievementStatus `json:"status" db:"status"`
	SubmittedAt        *time.Time        `json:"submittedAt" db:"submitted_at"`
	VerifiedAt         *time.Time        `json:"verifiedAt" db:"verified_at"`
	VerifiedBy         *uuid.UUID        `json:"verifiedBy" db:"verified_by"`
	RejectionNote      *string           `json:"rejectionNote" db:"rejection_note"`
	CreatedAt          time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time         `json:"updatedAt" db:"updated_at"`
}

type CreateAchievementReferenceRequest struct {
	StudentID          uuid.UUID `json:"studentId" validate:"required"`
	MongoAchievementID string    `json:"mongoAchievementId" validate:"required,len=24"`
}

type RejectAchievementRequest struct {
	RejectionNote string `json:"rejectionNote" validate:"required"`
}