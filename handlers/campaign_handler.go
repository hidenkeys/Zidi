package handlers

import (
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/utils"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type payment struct {
	Name         string
	Amount       string
	CampaignName string
	PaystackLink string
}

func (s Server) GenerateTokens(c *fiber.Ctx, id openapi_types.UUID) error {
	ctx := context.Background()

	campaign, err := s.campaignService.GetCampaignByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(api.Error{
			ErrorCode: "404",
			Message:   "Campaign not found",
		})
	}

	// Check if the campaign is active
	if campaign.Status != "active" {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Campaign is not active",
		})
	}

	// Generate the tokens
	tokens := utils.GenerateTokens(campaign.CharacterType, campaign.CouponLength, campaign.CouponNumber)

	// Store the tokens
	for _, token := range tokens {
		coupon := api.Coupon{
			CampaignId: id,
			Code:       token,
			Redeemed:   false,
		}
		_, err := s.campaignService.CreateCoupon(ctx, &coupon)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.Error{
				ErrorCode: "500",
				Message:   "Failed to store tokens",
			})
		}
	}

	// Respond with the generated tokens
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"campaignId": id,
		"tokens":     tokens,
	})
}

func (s Server) GetAllCampaigns(c *fiber.Ctx, params api.GetAllCampaignsParams) error {
	limit := 10
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}
	response, err := s.campaignService.GetAllCampaigns(context.Background(), limit, offset)
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
	//userClaims, ok := c.Locals("user").(middleware.UserClaims)
	//if !ok {
	//	return c.Status(http.StatusUnauthorized).JSON(api.Error{
	//		ErrorCode: "401",
	//		Message:   "Unauthorized - Invalid token",
	//	})
	//}
	limit := 10
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	response, err := s.campaignService.GetCampaignsByOrganization(context.Background(), params.OrganizationId, limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) GetCouponsByCampaign(c *fiber.Ctx, id openapi_types.UUID) error {
	response, err := s.campaignService.GetAllCoupons(context.Background(), id)
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

	pay := ""
	if campaign.Status == "pending" {
		response, err := s.orgService.GetOrganizationByID(context.Background(), campaign.OrganizationId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.Error{
				ErrorCode: "500",
				Message:   err.Error(),
			})
		}
		paymentLink, err := utils.CreatePaystackPaymentLink(string(response.Email), int(campaign.Amount)*100, id.String(), campaign.OrganizationId.String())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.Error{
				ErrorCode: "500",
				Message:   err.Error(),
			})
		}
		tmp := payment{
			Name:         response.ContactPersonName,
			Amount:       strconv.Itoa(int(campaign.Amount) * 100),
			CampaignName: campaign.CampaignName,
			PaystackLink: paymentLink,
		}

		tmpl, err := template.ParseFiles("Zidi-payment-email-template.html")
		if err != nil {
			log.Fatalf("Error loading template: %v", err)
		}

		// Parse the template with the receipt data
		var tpl bytes.Buffer
		if err := tmpl.Execute(&tpl, tmp); err != nil {
			log.Fatalf("Error executing template: %v", err)
		}

		// Convert parsed template to a string
		createBody := tpl.String()

		err = utils.SendEmail(string(response.Email), "Complete your "+campaign.CampaignName+" Campaign Payment", createBody)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.Error{
				ErrorCode: "500",
				Message:   err.Error(),
			})
		}
		pay = pay
	}

	response, err := s.campaignService.UpdateCampaign(context.Background(), id, campaign)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response, pay)
}
