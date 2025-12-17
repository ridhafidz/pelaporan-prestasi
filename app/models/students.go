package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"userId" db:"user_id"`
	StudentID      string     `json:"studentId" db:"student_id"`
	ProgramStudy   string     `json:"programStudy" db:"program_study"`
	AcademicYear   string     `json:"academicYear" db:"academic_year"`
	AdvisorID      *uuid.UUID `json:"advisorId" db:"advisor_id"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
	Name           string     `json:"name,omitempty" db:"name"`
	Email          string     `json:"email,omitempty" db:"email"`
	EnrollmentYear int        `json:"enrollmentYear,omitempty" db:"enrollment_year"`
	IsActive       bool       `json:"isActive,omitempty" db:"is_active"`
}

type CreateStudentRequest struct {
	UserID       uuid.UUID `json:"userId" validate:"required"`
	StudentID    string    `json:"studentId" validate:"required,max=20"`
	ProgramStudy string    `json:"programStudy" validate:"required,max=100"`
	AcademicYear string    `json:"academicYear" validate:"required,max=10"`
}

type UpdateAdvisorRequest struct {
	AdvisorID *uuid.UUID `json:"advisorId" db:"advisor_id"`
}
type StudentDetailResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	FullName     string    `json:"fullName"`
	StudentID    string    `json:"studentId"`
	ProgramStudy string    `json:"programStudy"`
	AcademicYear string    `json:"academicYear"`
	AdvisorName  string    `json:"advisorName"`
}

type StudentAchievementResponse struct {
	ID          string     `json:"achievement_id"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submittedAt,omitempty"`
	VerifiedAt  *time.Time `json:"verifiedAt,omitempty"`
}
