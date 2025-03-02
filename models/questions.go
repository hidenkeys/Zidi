package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CampaignID uuid.UUID `gorm:"not null;index"` // Links to Campaign
	Text       string    `gorm:"type:text;not null"`
	Type       string    `gorm:"not null; check:type IN ('text', 'multiple_choice', 'rating')"`
	Options    *string   `gorm:"type:jsonb"` // Stores multiple-choice options if applicable

	// Relations
	Responses []Response `gorm:"foreignKey:QuestionID"`
}
