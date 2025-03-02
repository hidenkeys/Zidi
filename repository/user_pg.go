package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
)

type UserRepoPG struct {
	db *gorm.DB
}

func NewUserRepoPG(db *gorm.DB) *UserRepoPG {
	return &UserRepoPG{db: db}
}

func (r *UserRepoPG) Create(user *models.User) (*models.User, error) {
	err := r.db.Create(user)
	if err.Error != nil {
		return nil, err.Error
	}
	return user, nil
}

func (r *UserRepoPG) GetByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("id = ?", id).First(user)
	if err.Error != nil {
		return nil, err.Error
	}
	return user, nil
}

func (r *UserRepoPG) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepoPG) UpdateByID(id uuid.UUID, user *models.User) (*models.User, error) {
	err := r.db.Model(&models.User{}).Where("id = ?", id).Updates(user)
	if err.Error != nil {
		return nil, err.Error
	}
	return user, nil
}

func (r *UserRepoPG) DeleteByID(id uuid.UUID) error {
	err := r.db.Where("id = ?", id).Delete(&models.User{})
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *UserRepoPG) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users)
	if err.Error != nil {
		return nil, err.Error
	}
	return users, nil
}

func (r *UserRepoPG) GetAllByOrganizationID(orgID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("organization_id = ?", orgID).Find(&users)
	if err.Error != nil {
		return nil, err.Error
	}
	return users, nil
}

func (r *UserRepoPG) UpdatePasswordByID(id uuid.UUID, password string) error {
	err := r.db.Model(&models.User{}).Where("id = ?", id).Update("password", password)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
