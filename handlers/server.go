package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/services"
	"github.com/hidenkeys/zidibackend/utils"
	"gorm.io/gorm"

	//openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	db              *gorm.DB
	orgService      *services.OrganizationService
	usrService      *services.UserService
	campaignService *services.CampaignService
	customerService *services.CustomerService
	questionService *services.QuestionService
	responseService *services.ResponseService
	paymentService  *services.PaymentService
}

type FlutterwaveWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		ID     int                    `json:"id"`
		Status string                 `json:"status"`
		TxRef  string                 `json:"tx_ref"`
		Amount int                    `json:"amount"`
		Meta   map[string]interface{} `json:"meta_data"` // Add this line
	} `json:"data"`
}

func (s Server) PostFlutterwaveWebhook(c *fiber.Ctx) error {
	// Read request body
	body := c.Body()

	fmt.Println(string(body))
	// Verify request signature
	signature := c.Get("verif-hash")
	secretHash := os.Getenv("FLW_SECRET_HASH") // Ensure this is set in your .env file

	if signature != secretHash {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid signature"})
	}

	// Parse JSON payload
	var payload FlutterwaveWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON payload"})
	}

	fmt.Println("3")
	// Log received webhook for debugging
	fmt.Printf("Received Flutterwave Webhook: %+v\n", payload)

	// Verify transaction ID
	transactionID := strconv.Itoa(payload.Data.ID)
	if transactionID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	// Verify transaction via Flutterwave API
	isVerified, err := utils.VerifyFlutterwaveTransaction(payload.Data.ID)
	if err != nil || !isVerified {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Transaction verification failed"})
	}

	fmt.Println(payload.Data)

	if payload.Data.Status == "successful" {
		fmt.Printf("Transaction %s verified. Processing payment...\n", transactionID)

		// Extract metadata from webhook
		//meta, ok := payload.Data.Meta.(map[string]interface{})
		//if !ok {
		//	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid metadata format"})
		//}

		// Extract campaign ID from metadata
		// Extract campaign ID from metadata
		campaignIDStr, exists := payload.Data.Meta["campaign_id"]
		if !exists {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing campaign ID"})
		}

		campaignID, err := uuid.Parse(campaignIDStr.(string))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid campaign ID"})
		}

		// Retrieve campaign details
		ctx := context.Background()
		campaign, err := s.campaignService.GetCampaignByID(ctx, campaignID)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Campaign not found"})
		}

		// Generate tokens
		tokens := utils.GenerateTokens(campaign.CharacterType, campaign.CouponLength, campaign.CouponNumber)

		// Store tokens in the database
		for _, token := range tokens {
			coupon := api.Coupon{
				CampaignId: campaignID,
				Code:       token,
				Redeemed:   false,
			}
			_, err := s.campaignService.CreateCoupon(ctx, &coupon)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store tokens"})
			}
		}

		// Save payment record
		payment := api.Payment{
			TransactionId:  transactionID,
			CampaignId:     campaignID,
			Amount:         float32(payload.Data.Amount),
			Status:         "successful",
			OrganizationId: campaign.OrganizationId,
			TransactionRef: payload.Data.TxRef,
		}

		if _, err := s.paymentService.CreatePayment(ctx, &payment); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save payment record"})
		}
		// Update campaign status to "active"
		campaign.Status = "active"
		if err, _ := s.campaignService.UpdateCampaign(ctx, campaignID, campaign); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update campaign status"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "Payment processed successfully, campaign activated, tokens generated",
			"tokens":  tokens,
		})
	}

	fmt.Printf("Transaction %s failed or pending. Skipping...\n", transactionID)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Webhook received"})
}

func NewServer(db *gorm.DB, orgService *services.OrganizationService, usrService *services.UserService, campaignService *services.CampaignService, customerService *services.CustomerService, questionService *services.QuestionService, responseService *services.ResponseService, paymentService *services.PaymentService) *Server {
	return &Server{
		db:              db,
		orgService:      orgService,
		usrService:      usrService,
		campaignService: campaignService,
		customerService: customerService,
		questionService: questionService,
		responseService: responseService,
		paymentService:  paymentService,
	}
}
