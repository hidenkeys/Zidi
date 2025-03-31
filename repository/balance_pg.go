package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
)

type BalanceRepoPG struct {
	db *gorm.DB
}

func NewBalanceRepoPG(db *gorm.DB) *BalanceRepoPG {
	return &BalanceRepoPG{db: db}
}

func (r *BalanceRepoPG) CreateBalance(balance *models.Balance) error {
	return r.db.Create(balance).Error
}

func (r *BalanceRepoPG) GetBalanceByCampaign(campaignId uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	err := r.db.Where("campaign_id = ?", campaignId).First(&balance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &balance, err
}

func (r *BalanceRepoPG) UpdateBalance(campaignId uuid.UUID, amount float64) error {
	return r.db.Model(&models.Balance{}).
		Where("campaign_id = ?", campaignId).
		Update("amount", amount).Error
}

func (r *BalanceRepoPG) GetAllBalances(limit, offset int) ([]models.Balance, error) {
	var balances []models.Balance
	err := r.db.Limit(limit).Offset(offset).Find(&balances).Error
	return balances, err
}
