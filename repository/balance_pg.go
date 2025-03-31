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

func (r *BalanceRepoPG) CreateBalance(balance *models.Balance) (*models.Balance, error) {
	if err := r.db.Create(balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepoPG) GetBalanceByCampaign(campaignId uuid.UUID) ([]models.Balance, error) {
	var balances []models.Balance
	err := r.db.Where("campaign_id = ?", campaignId).Find(&balances).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return balances, err
}

func (r *BalanceRepoPG) UpdateBalance(campaignId uuid.UUID, balance *models.Balance) (*models.Balance, error) {
	if err := r.db.Model(&models.Balance{}).
		Where("campaign_id = ?", campaignId).
		Updates(balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepoPG) GetAllBalances() ([]models.Balance, error) {
	var balances []models.Balance
	err := r.db.Find(&balances).Error
	return balances, err
}
