package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
)

type OrganizationRepository interface {
	Create(org *models.Organization) (*models.Organization, error)
	GetByID(id uuid.UUID) (*models.Organization, error)
	GetAllById() ([]models.Organization, error)
	GetByName(name string) ([]models.Organization, error)
	UpdateByID(id uuid.UUID, org *models.Organization) (*models.Organization, error)
	DeleteByID(id uuid.UUID) error
}
