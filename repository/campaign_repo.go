package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
)

type CampaignRepository interface {
	Create(campaign *models.Campaign) (*models.Campaign, error)
	GetAll() ([]models.Campaign, error)
	GetAllByOrganization(orgID uuid.UUID) ([]models.Campaign, error)
	GetCouponByCampaign(campaignId uuid.UUID) ([]models.Coupon, error)
	GetByID(id uuid.UUID) (*models.Campaign, error)
	UpdateByID(id uuid.UUID, campaign *models.Campaign) (*models.Campaign, error)
	DeleteByID(id uuid.UUID) error
	CreateCoupon(coupon *models.Coupon) (*models.Coupon, error)
}
