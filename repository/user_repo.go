package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdateByID(id uuid.UUID, user *models.User) (*models.User, error)
	DeleteByID(id uuid.UUID) error
	UpdatePasswordByID(id uuid.UUID, password string) error
	GetAllByOrganizationID(orgID uuid.UUID) ([]models.User, error)
}
