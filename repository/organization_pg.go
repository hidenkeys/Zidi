package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
)

type OrganizationRepoPG struct {
	db *gorm.DB
}

func NewOrganizationRepoPG(db *gorm.DB) *OrganizationRepoPG {
	return &OrganizationRepoPG{db: db}
}

func (r *OrganizationRepoPG) Create(org *models.Organization) (*models.Organization, error) {
	if err := r.db.Create(org).Error; err != nil {
		return nil, err
	}
	return org, nil
}

func (r *OrganizationRepoPG) GetByID(id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.Where("id = ?", id).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepoPG) GetAllById() ([]models.Organization, error) {
	var orgs []models.Organization
	if err := r.db.Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *OrganizationRepoPG) GetByName(name string) ([]models.Organization, error) {
	var orgs []models.Organization
	if err := r.db.Where("company_name ILIKE ?", "%"+name+"%").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *OrganizationRepoPG) UpdateByID(id uuid.UUID, org *models.Organization) (*models.Organization, error) {
	if err := r.db.Model(&models.Organization{}).Where("id = ?", id).Updates(org).Error; err != nil {
		return nil, err
	}
	return org, nil
}

func (r *OrganizationRepoPG) DeleteByID(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Organization{}).Error; err != nil {
		return err
	}
	return nil
}
