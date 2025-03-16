package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName      string
	LastName       string
	Phone          string
	Email          string
	Feedback       string
	Network        string
	Amount         float64
	Status         string
	OrganizationID uuid.UUID `gorm:"constraint:OnDelete:CASCADE;"`
	CampaignID     uuid.UUID `gorm:"constraint:OnDelete:CASCADE;"`
}
