package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
)

type CustomerRepository interface {
	Create(customer *models.Customer) (*models.Customer, error)
	UpdateByID(id uuid.UUID, customer *models.Customer) (*models.Customer, error)
	DeleteByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*models.Customer, error)
	GetAll() ([]models.Customer, error)
	GetAllByOrganization(orgID uuid.UUID) ([]models.Customer, error)
}
