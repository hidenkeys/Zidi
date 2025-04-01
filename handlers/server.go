package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/services"
	"github.com/hidenkeys/zidibackend/utils"
	"gorm.io/gorm"
	"html/template"
	"log"

	//openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"os"
	"strconv"
)

type ppayment struct {
	Name            string
	Amount          string
	CampaignName    string
	TelegramBotLink string
}

type Server struct {
	db                 *gorm.DB
	orgService         *services.OrganizationService
	usrService         *services.UserService
	campaignService    *services.CampaignService
	customerService    *services.CustomerService
	questionService    *services.QuestionService
	responseService    *services.ResponseService
	paymentService     *services.PaymentService
	transactionService *services.TransactionService
	balanceService     *services.BalanceService
}

type FlutterwaveWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		ID                int    `json:"id"`
		TxRef             string `json:"tx_ref"`
		FlwRef            string `json:"flw_ref"`
		Amount            int    `json:"amount"`
		Currency          string `json:"currency"`
		ChargedAmount     int    `json:"charged_amount"`
		AppFee            int    `json:"app_fee"`
		MerchantFee       int    `json:"merchant_fee"`
		ProcessorResponse string `json:"processor_response"`
		AuthModel         string `json:"auth_model"`
		IP                string `json:"ip"`
		Narration         string `json:"narration"`
		Status            string `json:"status"`
		PaymentType       string `json:"payment_type"`
		CreatedAt         string `json:"created_at"`
		AccountID         int    `json:"account_id"`
		Customer          struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Email       string `json:"email"`
			CreatedAt   string `json:"created_at"`
		} `json:"customer"`
	} `json:"data"`
	Meta struct {
		CheckoutInitAddress     string `json:"__CheckoutInitAddress"`
		CampaignID              string `json:"campaign_id"`
		OrganizationID          string `json:"organization_id"`
		OriginatorAccountNumber string `json:"originatoraccountnumber"`
		OriginatorName          string `json:"originatorname"`
		BankName                string `json:"bankname"`
		OriginatorAmount        string `json:"originatoramount"`
	} `json:"meta_data"`
}

func (s Server) PostFlutterwaveWebhook(c *fiber.Ctx) error {
	// Read request body
	body := c.Body()

	// Convert body to a string and log it
	bodyStr := string(body)
	fmt.Printf("Request Body as Text: %s\n", bodyStr)

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
	// Log received webhook for debugging
	fmt.Printf("Received Flutterwave Webhook: %+v\n", payload)

	fmt.Println("1")
	// Verify transaction ID
	transactionID := strconv.Itoa(payload.Data.ID)
	if transactionID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	fmt.Println("2")
	// Verify transaction via Flutterwave API
	isVerified, err := utils.VerifyFlutterwaveTransaction(payload.Data.ID)
	if err != nil || !isVerified {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Transaction verification failed"})
	}

	fmt.Println("3")
	// Extract campaign ID from metadata
	campaignIDStr := payload.Meta.CampaignID
	fmt.Println("campaig id", campaignIDStr)
	if campaignIDStr == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing campaign ID"})
	}

	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid campaign ID"})
	}

	// Retrieve campaign details
	ctx := context.Background()

	campaign, err := s.campaignService.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Campaign not found"})
	}

	fmt.Println("4")

	fmt.Println("payload amount ", payload.Data.Amount)
	fmt.Println("campaig amunt", campaign.Price)

	if payload.Data.Status == "successful" && float32(payload.Data.Amount) == campaign.Price {
		fmt.Printf("Transaction %s verified. Processing payment...\n", transactionID)

		fmt.Println("5")
		// Extract metadata from webhook
		//meta, ok := payload.Data.Meta.(map[string]interface{})
		//if !ok {
		//	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid metadata format"})
		//}
		// Create balance record

		balance := api.Balance{
			CampaignId:      campaignID,
			StartingBalance: float32(payload.Data.Amount),
			Amount:          float32(payload.Data.Amount),
		}
		if _, err := s.balanceService.CreateBalance(ctx, &balance); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create balance record"})
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

		response, err := s.orgService.GetOrganizationByID(context.Background(), campaign.OrganizationId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.Error{
				ErrorCode: "500",
				Message:   err.Error(),
			})
		}

		link := "https://t.me/zidipromobot?start=" + campaignID.String()
		tmp := ppayment{
			Name:            response.ContactPersonName,
			Amount:          strconv.Itoa(int(campaign.Amount) * 100),
			CampaignName:    campaign.CampaignName,
			TelegramBotLink: link, // Updated to use Flutterwave link
		}

		tmpl, err := template.ParseFiles("Zidi-payment-successful-email-template.html")
		if err != nil {
			log.Fatalf("Error loading template: %v", err)
		}

		var tpl bytes.Buffer
		if err := tmpl.Execute(&tpl, tmp); err != nil {
			log.Fatalf("Error executing template: %v", err)
		}

		createBody := tpl.String()

		err = utils.SendEmail(string(response.Email), "Complete your "+campaign.CampaignName+" Campaign Payment", createBody)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(api.Error{
				ErrorCode: "500",
				Message:   err.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "Payment processed successfully, campaign activated, tokens generated",
			"tokens":  tokens,
		})
	}

	fmt.Printf("Transaction %s failed or pending. Skipping...\n", transactionID)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Webhook received"})
}

func NewServer(db *gorm.DB, balanceService *services.BalanceService, transactionService *services.TransactionService, orgService *services.OrganizationService, usrService *services.UserService, campaignService *services.CampaignService, customerService *services.CustomerService, questionService *services.QuestionService, responseService *services.ResponseService, paymentService *services.PaymentService) *Server {
	return &Server{
		db:                 db,
		balanceService:     balanceService,
		transactionService: transactionService,
		orgService:         orgService,
		usrService:         usrService,
		campaignService:    campaignService,
		customerService:    customerService,
		questionService:    questionService,
		responseService:    responseService,
		paymentService:     paymentService,
	}
}
