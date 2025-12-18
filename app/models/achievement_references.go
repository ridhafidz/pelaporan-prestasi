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

type AchievementTypeStat struct {
	AchievementType string `json:"achievementType"`
	Total           int    `json:"total"`
}

type AchievementPeriodStat struct {
	Period string `json:"period"` // contoh: 2024-01
	Total  int    `json:"total"`
}

type TopStudentStat struct {
	StudentID  uuid.UUID `json:"studentId"`
	FullName   string    `json:"fullName"`
	TotalPoint float64   `json:"totalPoint"`
}

type CompetitionLevelStat struct {
	Level string `json:"level"`
	Total int    `json:"total"`
}

type AchievementStatisticsResponse struct {
	ByType            []AchievementTypeStat   `json:"byType"`
	ByPeriod          []AchievementPeriodStat `json:"byPeriod"`
	TopStudents       []TopStudentStat        `json:"topStudents"`
	CompetitionLevels []CompetitionLevelStat  `json:"competitionLevels"`
}

type StudentReportResponse struct {
	StudentID  uuid.UUID `json:"studentId"`
	TotalPoint float64   `json:"totalPoint"`
}
