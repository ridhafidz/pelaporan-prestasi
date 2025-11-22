package models

import "time"

type AchievementReference struct {
	ID               string    `json:"id" db:"id"`
	StudentID        string    `json:"student_id" db:"student_id"`
	MongoAchievementID string  `json:"mongo_achievement_id" db:"mongo_achievement_id"`
	Status           string    `json:"status" db:"status"` // draft, submitted, verified, rejected
	SubmittedAt      time.Time `json:"submitted_at" db:"submitted_at"`
	VerifiedAt       time.Time `json:"verified_at" db:"verified_at"`
	VerifiedBy       string    `json:"verified_by" db:"verified_by"`
	RejectionNote    string    `json:"rejection_note" db:"rejection_note"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
