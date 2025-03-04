package models

import (
	"github.com/google/uuid"
)

type Campaign struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CampaignName   string    `gorm:"not null"`
	CouponID       string
	CharacterType  string
	CouponLength   int
	CouponNumber   int
	OrganizationID uuid.UUID
	WelcomeMessage string `gorm:"type:text"`
	QuestionNumber int
	Amount         float64
	Status         string

	// Relations
	Questions []Question `gorm:"foreignKey:CampaignID"`
}
