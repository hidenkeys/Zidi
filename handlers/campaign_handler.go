package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/hidenkeys/zidibackend/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s Server) GetAllCampaigns(c *fiber.Ctx) error {
	response, err := s.campaignService.GetAllCampaigns(context.Background())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) CreateCampaign(c *fiber.Ctx) error {
	campaign := new(api.Campaign)
	if err := c.BodyParser(campaign); err != nil {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid request body",
		})
	}

	response, err := s.campaignService.CreateCampaign(context.Background(), *campaign)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusCreated).JSON(response)
}

func (s Server) GetCampaignsByOrganization(c *fiber.Ctx, params api.GetCampaignsByOrganizationParams) error {
	response, err := s.campaignService.GetCampaignsByOrganization(context.Background(), params.OrganizationId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) DeleteCampaign(c *fiber.Ctx, id openapi_types.UUID) error {
	err := s.campaignService.DeleteCampaign(context.Background(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusNoContent).JSON(nil)
}

func (s Server) GetCampaignById(c *fiber.Ctx, id openapi_types.UUID) error {
	response, err := s.campaignService.GetCampaignByID(context.Background(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) UpdateCampaign(c *fiber.Ctx, id openapi_types.UUID) error {
	campaign := new(api.Campaign)
	if err := c.BodyParser(campaign); err != nil {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid request body",
		})
	}

	response, err := s.campaignService.UpdateCampaign(context.Background(), id, campaign)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}
