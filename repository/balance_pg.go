package repository

import (
	"context"
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

func (r *BalanceRepoPG) CreateBalance(ctx context.Context, balance *models.Balance) error {
	return r.db.WithContext(ctx).Create(balance).Error
}

func (r *BalanceRepoPG) GetBalanceByCampaign(ctx context.Context, campaignId uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignId).First(&balance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &balance, err
}

func (r *BalanceRepoPG) UpdateBalance(ctx context.Context, campaignId uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&models.Balance{}).
		Where("campaign_id = ?", campaignId).
		Update("amount", amount).Error
}

func (r *BalanceRepoPG) GetAllBalances(ctx context.Context, limit, offset int) ([]models.Balance, error) {
	var balances []models.Balance
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&balances).Error
	return balances, err
}
