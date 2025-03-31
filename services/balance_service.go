package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/models"
	"github.com/hidenkeys/zidibackend/repository"
)

type BalanceService struct {
	balanceRepo repository.BalanceRepository
}

func NewBalanceService(balanceRepo repository.BalanceRepository) *BalanceService {
	return &BalanceService{balanceRepo: balanceRepo}
}

func (s *BalanceService) CreateBalance(ctx context.Context, req *api.Balance) (*api.Balance, error) {
	newBalance := &models.Balance{
		CampaignId:      req.CampaignId,
		StartingBalance: float64(req.StartingBalance),
		Amount:          float64(req.Amount),
	}

	createdBalance, err := s.balanceRepo.Create(newBalance)
	if err != nil {
		return nil, err
	}

	return mapToAPIBalance(createdBalance), nil
}

func (s *BalanceService) GetBalanceByCampaign(ctx context.Context, campaignId uuid.UUID) (*api.Balance, error) {
	balance, err := s.balanceRepo.GetByCampaignID(campaignId)
	if err != nil {
		return nil, err
	}
	if balance == nil {
		return nil, nil
	}

	return mapToAPIBalance(balance), nil
}

func (s *BalanceService) GetAllBalances(ctx context.Context, limit, offset int) ([]api.Balance, error) {
	balances, err := s.balanceRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var finalBalances []api.Balance
	for _, balance := range balances {
		finalBalances = append(finalBalances, *mapToAPIBalance(&balance))
	}

	return finalBalances, nil
}

func (s *BalanceService) UpdateBalance(ctx context.Context, campaignId uuid.UUID, amount float64) (*api.Balance, error) {
	updatedBalance, err := s.balanceRepo.UpdateByCampaignID(campaignId, &models.Balance{Amount: amount})
	if err != nil {
		return nil, err
	}

	//updatedBalance, err := s.balanceRepo.GetByCampaignID(campaignId)
	//if err != nil {
	//	return nil, err
	//}

	return mapToAPIBalance(updatedBalance), nil
}

func mapToAPIBalance(balance *models.Balance) *api.Balance {
	return &api.Balance{
		Id:              balance.ID,
		CampaignId:      balance.CampaignId,
		StartingBalance: float32(balance.StartingBalance),
		Amount:          float32(balance.Amount),
	}
}
