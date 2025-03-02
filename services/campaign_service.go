package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/models"
	"github.com/hidenkeys/zidibackend/repository"
)

type CampaignService struct {
	campaignRepo repository.CampaignRepository
}

func NewCampaignService(campaignRepo repository.CampaignRepository) *CampaignService {
	return &CampaignService{campaignRepo: campaignRepo}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, req api.Campaign) (*api.Campaign, error) {
	newCampaign := &models.Campaign{
		ID:             req.Id,
		CampaignName:   req.CampaignName,
		CouponID:       req.CouponId,
		CharacterType:  req.CharacterType,
		CouponLength:   req.CouponLength,
		OrganizationID: req.OrganizationId,
		WelcomeMessage: req.WelcomeMessage,
		QuestionNumber: req.QuestionNumber,
		Amount:         float64(req.Amount),
		Status:         req.Status,
	}

	campaign, err := s.campaignRepo.Create(newCampaign)
	if err != nil {
		return nil, err
	}

	return mapToAPICampaign(campaign), nil
}

func (s *CampaignService) GetAllCampaigns(ctx context.Context) ([]api.Campaign, error) {
	campaigns, err := s.campaignRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var finalCampaigns []api.Campaign
	for _, campaign := range campaigns {
		finalCampaigns = append(finalCampaigns, *mapToAPICampaign(&campaign))
	}

	return finalCampaigns, nil
}

func (s *CampaignService) GetCampaignsByOrganization(ctx context.Context, orgID uuid.UUID) ([]api.Campaign, error) {
	campaigns, err := s.campaignRepo.GetAllByOrganization(orgID)
	if err != nil {
		return nil, err
	}

	var finalCampaigns []api.Campaign
	for _, campaign := range campaigns {
		finalCampaigns = append(finalCampaigns, *mapToAPICampaign(&campaign))
	}

	return finalCampaigns, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, id uuid.UUID) (*api.Campaign, error) {
	campaign, err := s.campaignRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return mapToAPICampaign(campaign), nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, id uuid.UUID, req *api.Campaign) (*api.Campaign, error) {
	updateData := &models.Campaign{
		CampaignName:   req.CampaignName,
		CouponID:       req.CouponId,
		CharacterType:  req.CharacterType,
		CouponLength:   req.CouponLength,
		OrganizationID: req.OrganizationId,
		WelcomeMessage: req.WelcomeMessage,
		QuestionNumber: req.QuestionNumber,
		Amount:         float64(req.Amount),
		Status:         req.Status,
	}

	updatedCampaign, err := s.campaignRepo.UpdateByID(id, updateData)
	if err != nil {
		return nil, err
	}

	return mapToAPICampaign(updatedCampaign), nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, id uuid.UUID) error {
	return s.campaignRepo.DeleteByID(id)
}

// Helper function to convert models.Campaign to api.Campaign
func mapToAPICampaign(campaign *models.Campaign) *api.Campaign {
	return &api.Campaign{
		Id:             campaign.ID,
		CampaignName:   campaign.CampaignName,
		CouponId:       campaign.CouponID,
		CharacterType:  campaign.CharacterType,
		CouponLength:   campaign.CouponLength,
		OrganizationId: campaign.OrganizationID,
		WelcomeMessage: campaign.WelcomeMessage,
		QuestionNumber: campaign.QuestionNumber,
		Amount:         float32(campaign.Amount),
		Status:         campaign.Status,
	}
}
