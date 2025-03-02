package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
)

type CustomerPG struct {
	db *gorm.DB
}

func NewCustomerRepoPG(db *gorm.DB) *CustomerPG {
	return &CustomerPG{db: db}
}

func (r *CustomerPG) Create(customer *models.Customer) (*models.Customer, error) {
	if err := r.db.Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *CustomerPG) UpdateByID(id uuid.UUID, customer *models.Customer) (*models.Customer, error) {
	if err := r.db.Model(&models.Customer{}).Updates(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *CustomerPG) DeleteByID(id uuid.UUID) error {
	if err := r.db.Delete(&models.Customer{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *CustomerPG) GetByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerPG) GetAll() ([]models.Customer, error) {
	var customers []models.Customer
	if err := r.db.Find(&customers).Error; err != nil {
		return nil, err
	}
	if customers == nil {
		customers = []models.Customer{}
	}
	return customers, nil
}

func (r *CustomerPG) GetAllByOrganization(orgID uuid.UUID) ([]models.Customer, error) {
	var customers []models.Customer
	if err := r.db.Where("organization_id = ?", orgID).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
