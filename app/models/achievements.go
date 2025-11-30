package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	FileName   string    `json:"fileName" bson:"fileName"`
	FileURL    string    `json:"fileURL" bson:"fileURL"`
	FileType   string    `json:"fileType" bson:"fileType"`
	UploadedAt time.Time `json:"uploadedAt" bson:"uploadedAt"`
}

type Period struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

type DynamicDetails struct {
	CompetitionName  string        `json:"competitionName,omitempty" bson:"competitionName,omitempty"`
	CompetitionLevel string        `json:"competitionLevel,omitempty" bson:"competitionLevel,omitempty"`
	Rank             int           `json:"rank,omitempty" bson:"rank,omitempty"`
	MedalType        string        `json:"medalType,omitempty" bson:"medalType,omitempty"`
	PublicationType  string        `json:"publicationType,omitempty" bson:"publicationType,omitempty"`
	PublicationTitle string        `json:"publicationTitle,omitempty" bson:"publicationTitle,omitempty"`
	Authors          []string      `json:"authors,omitempty" bson:"authors,omitempty"`
	Publisher        string        `json:"publisher,omitempty" bson:"publisher,omitempty"`
	ISSN             string        `json:"issn,omitempty" bson:"issn,omitempty"`
	OrganizationName string        `json:"organizationName,omitempty" bson:"organizationName,omitempty"`
	Position         string        `json:"position,omitempty" bson:"position,omitempty"`
	Period           Period        `json:"period,omitempty" bson:"period,omitempty"`
	CertificationName   string        `json:"certificationName,omitempty" bson:"certificationName,omitempty"`
	IssuedBy            string        `json:"issuedBy,omitempty" bson:"issuedBy,omitempty"`
	CertificationNumber string        `json:"certificationNumber,omitempty" bson:"certificationNumber,omitempty"`
	ValidUntil       time.Time     `json:"validUntil,omitempty" bson:"validUntil,omitempty"`
	EventDate        time.Time     `json:"eventDate,omitempty" bson:"eventDate,omitempty"`
	Location         string        `json:"location,omitempty" bson:"location,omitempty"`
	Organizer        string        `json:"organizer,omitempty" bson:"organizer,omitempty"`
	Score            float64       `json:"score,omitempty" bson:"score,omitempty"`
	CustomFields     primitive.M   `json:"customFields,omitempty" bson:"customFields,omitempty"`
}

type Achievement struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StudentID       string             `json:"studentId" bson:"studentId"`
	AchievementType string             `json:"achievementType" bson:"achievementType"`
	Title           string             `json:"title" bson:"title"`
	Description     string             `json:"description" bson:"description"`
	Details         DynamicDetails     `json:"details" bson:"details"`
	Attachments     []Attachment       `json:"attachments" bson:"attachments"`
	Tags            []string           `json:"tags" bson:"tags"`
	Points          float64            `json:"points" bson:"points"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt       *time.Time         `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
