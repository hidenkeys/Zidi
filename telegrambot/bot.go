package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hidenkeys/zidibackend/utils"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	//"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type createOrgaization struct {
	Name       string
	CouponCode string
}
type Session struct {
	CampaignID      uuid.UUID
	OrganizationID  uuid.UUID
	Amount          float64
	Customer        models.Customer
	Step            int
	Questions       []models.Question
	CurrentQuestion int
	Responses       []models.Response
}

var sessions = make(map[int64]*Session)

// StartBot initializes and runs the Telegram bot
// StartBot initializes and runs the Telegram bot with a database connection
func StartBot(db *gorm.DB) {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal("❌ Error loading .env file")
	//}

	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_API_KEY"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("❌ Error initializing bot: ", err)
	}

	log.Println("🚀 Telegram Bot is running...")

	bot.Handle("/start", func(c tele.Context) error {
		args := c.Args()
		if len(args) == 0 {
			return c.Send("❌ Please provide a valid campaign ID.")
		}

		campaignID, err := uuid.Parse(args[0])
		if err != nil {
			return c.Send("❌ Invalid campaign ID format.")
		}

		var campaign models.Campaign
		if err := db.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Send("❌ Campaign not found.")
			}
			return c.Send("❌ An error occurred while fetching the campaign.")
		}

		var questions []models.Question
		if err := db.Where("campaign_id = ?", campaignID).Find(&questions).Error; err != nil {
			return c.Send("❌ Failed to fetch questions.")
		}

		user := c.Sender()

		// Check if the customer already exists for the same campaign
		var existingCustomer models.Customer
		if err := db.Where("phone = ? AND campaign_id = ?", user.Recipient(), campaign.ID).First(&existingCustomer).Error; err == nil {
			return c.Send("✅ You have already registered for this campaign. Stay tuned for the next one!")
		}

		sessions[user.ID] = &Session{
			CampaignID:     campaign.ID,
			OrganizationID: campaign.OrganizationID,
			Amount:         campaign.Amount,
			Customer: models.Customer{
				Status:     "inactive",
				CampaignID: campaign.ID,
			},
			Step:            1,
			Questions:       questions,
			CurrentQuestion: 0,
			Responses:       []models.Response{},
		}

		return c.Send(fmt.Sprintf("👋 Welcome!\n%s\n\nLet's get started by gathering your details.\nWhat's your first name?", campaign.WelcomeMessage))
	})

	// Pass the db to handleResponses
	bot.Handle(tele.OnText, func(c tele.Context) error {
		return handleResponses(c, db)
	})

	bot.Start()
}
func parseOptions(optionsJSON []byte) ([]string, error) {
	var options []string
	if err := json.Unmarshal(optionsJSON, &options); err != nil {
		return nil, err
	}
	return options, nil
}

// Helper to create option buttons
func createOptionButtons(options []string) [][]tele.ReplyButton {
	var buttons [][]tele.ReplyButton
	for _, option := range options {
		btn := tele.ReplyButton{Text: option}
		buttons = append(buttons, []tele.ReplyButton{btn}) // One button per row
	}
	return buttons
}

func handleResponses(c tele.Context, db *gorm.DB) error {
	userID := c.Sender().ID
	session, exists := sessions[userID]
	if !exists {
		return c.Send("❌ Please start with /start and a valid campaign ID.")
	}

	switch session.Step {
	case 1:
		session.Customer.FirstName = c.Text()
		session.Step++
		return c.Send("📛 What's your last name?")

	case 2:
		session.Customer.LastName = c.Text()
		session.Step++
		return c.Send("📧 What’s your email address?")

	case 3:
		session.Customer.Email = c.Text()
		session.Step++
		var coupon models.Coupon

		// Find the first unredeemed coupon for a specific campaign
		if err := db.
			Where("campaign_id = ? AND redeemed = false", session.CampaignID).
			First(&coupon).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				return c.Send("❌ No available coupons at the moment.")
			}

			log.Println("❌ Error retrieving coupon:", err)
			return c.Send("❌ An error occurred while fetching a coupon. Please try again later.")
		}

		// Coupon found, send it to the user
		tmp := createOrgaization{
			Name:       session.Customer.FirstName + " " + session.Customer.LastName,
			CouponCode: coupon.Code,
		}

		tmpl, err := template.ParseFiles("Zidi-coupon-code-email-template.html")
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

		err = utils.SendEmail0(string(session.Customer.Email), "Your Zidi Campaign Coupon Code", createBody)
		if err != nil {
			return c.Send("❌ An error occurred while fetching a coupon. Please try again later.")
		}
		return c.Send("📞 Please provide your phone number:")

	case 4:
		session.Customer.Phone = c.Text()
		session.Step++

		networkKeyboard := &tele.ReplyMarkup{ResizeKeyboard: true}
		btnMTN := networkKeyboard.Text("mtn")
		btnGlo := networkKeyboard.Text("glo")
		btnAirtel := networkKeyboard.Text("airtel")
		btnEtisalat := networkKeyboard.Text("etisalat")
		networkKeyboard.Reply(networkKeyboard.Row(btnMTN, btnGlo), networkKeyboard.Row(btnAirtel, btnEtisalat))

		return c.Send("📶 Which network provider do you use?", networkKeyboard)

	case 5:
		network := strings.ToLower(c.Text())
		validNetworks := map[string]bool{
			"mtn":      true,
			"glo":      true,
			"airtel":   true,
			"etisalat": true,
		}

		if !validNetworks[network] {
			return c.Send("❌ Invalid network. Please select from the options provided.")
		}

		session.Customer.Network = strings.ToUpper(network) // optional: normalize
		session.Step++
		return c.Send("💬 How do you like the experience so far?", &tele.ReplyMarkup{RemoveKeyboard: true})

	case 6:
		session.Customer.Feedback = c.Text()
		session.Customer.OrganizationID = session.OrganizationID
		session.Customer.Amount = session.Amount

		// Save customer to DB
		if err := saveCustomer(db, &session.Customer); err != nil {
			log.Println("❌ Error saving customer:", err)
			return c.Send("❌ customer already exists for this campaign.")
		}

		// Move to questions if available
		if len(session.Questions) > 0 {
			session.Step = 7
			firstQuestion := session.Questions[0]

			// Parse options if present
			options, err := parseOptions(firstQuestion.Options)
			if err != nil {
				log.Println("❌ Error parsing options:", err)
				return c.Send("❌ Error retrieving question options. Please try again later.")
			}

			// Send question with buttons if there are options
			if len(options) > 0 && firstQuestion.Type == "multiple_choice" {
				btns := createOptionButtons(options)
				return c.Send(fmt.Sprintf("📋 %s", firstQuestion.Text), &tele.ReplyMarkup{ReplyKeyboard: btns})
			}

			// Otherwise, just send the question
			return c.Send(fmt.Sprintf("📋 %s", firstQuestion.Text))
		}

		// If no questions, finalize
		session.Step = 8
		return c.Send("🎟 Please enter your coupon code to complete the registration.")

	case 7:
		// Capture the user's response
		question := session.Questions[session.CurrentQuestion]
		session.Responses = append(session.Responses, models.Response{
			CustomerID: session.Customer.ID,
			QuestionID: question.ID,
			Answer:     c.Text(),
		})

		// If there are more questions, move to the next one
		if session.CurrentQuestion+1 < len(session.Questions) {
			session.CurrentQuestion++
			nextQuestion := session.Questions[session.CurrentQuestion]

			// Parse and show options if applicable
			options, err := parseOptions(nextQuestion.Options)
			if err != nil {
				log.Println("❌ Error parsing options:", err)
				return c.Send("❌ Error retrieving question options. Please try again later.")
			}

			if len(options) > 0 && nextQuestion.Type == "multiple_choice" {
				btns := createOptionButtons(options)
				return c.Send(nextQuestion.Text, &tele.ReplyMarkup{ReplyKeyboard: btns})
			}

			return c.Send(nextQuestion.Text)
		}

		// Save responses to DB
		if err := saveResponses(db, session.Responses); err != nil {
			log.Println("❌ Error saving responses:", err)
			return c.Send("❌ An error occurred while saving your responses.")
		}

		// Move to coupon validation step
		session.Step = 8
		clearKeyboard := &tele.ReplyMarkup{RemoveKeyboard: true}
		//return c.Send("🎟 Please enter your coupon code to complete the registration.", clearKeyboard)
		return c.Send("🎟 Please enter the coupon code sent to your email.", clearKeyboard)

	case 8: // Step 9: Validate coupon code
		couponCode := c.Text()

		var coupon models.Coupon
		if err := db.Where("code = ? AND campaign_id = ?", couponCode, session.CampaignID).First(&coupon).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Send("❌ Invalid coupon code. Please try again.")
			}
			log.Println("❌ Error retrieving coupon:", err)
			return c.Send("❌ An error occurred while checking the coupon. Please try again later.")
		}

		// Check if the coupon has already been redeemed
		if coupon.Redeemed {
			return c.Send("❌ This coupon has already been redeemed. Please check and try again.")
		}

		// Mark coupon as redeemed
		now := time.Now()
		coupon.Redeemed = true
		coupon.RedeemedAt = &now
		if err := db.Save(&coupon).Error; err != nil {
			log.Println("❌ Error updating coupon:", err)
			return c.Send("❌ An error occurred while redeeming your coupon. Please try again later.")
		}

		airtimeRespose, err := utils.SendAirtime(fmt.Sprintf("%.0f", session.Amount), session.Customer.Network, session.Customer.Phone)
		if err != nil {
			log.Println("❌ Error sending airtime:", err)
			return c.Send("❌ An error occurred while sending your airtime. Please try again later.")
		}

		commissionFloat, err := strconv.ParseFloat(airtimeRespose.Commission, 32)
		if err != nil {
			fmt.Println("error converting commission:", err)
			return nil
		}
		//// Create a new transaction after sending airtime
		//transactionInput := &api.TransactionInput{
		//	OrganizationId: session.OrganizationID,
		//	CampaignId:     session.CampaignID,
		//	CustomerId:     session.Customer.ID,
		//	Network:        airtimeRespose.Network,
		//	PhoneNumber:    airtimeRespose.Phone,
		//	TxReference:    airtimeRespose.RequestID,
		//	Amount:         float32(session.Amount), // Amount is already in float64
		//	Type:           "airtime",
		//	Commisson:      float32(commissionFloat),
		//}

		// Create transaction directly
		tx := models.Transaction{
			OrganizationID: session.OrganizationID,
			CampaignID:     session.CampaignID,
			CustomerID:     session.Customer.ID,
			Network:        airtimeRespose.Network,
			PhoneNumber:    airtimeRespose.Phone,
			TxReference:    airtimeRespose.RequestID,
			Status:         airtimeRespose.Status,
			Amount:         session.Amount,
			Type:           "airtime",
			Commisson:      commissionFloat,
		}

		if err := db.Create(&tx).Error; err != nil {
			log.Println("❌ Error creating transaction directly:", err)
			return c.Send("❌ An error occurred while processing your transaction. Please try again later.")
		}

		//transactionRepo := repository.NewTransactionRepoPG(db)
		//transactionService := services.NewTransactionService(transactionRepo) // assuming you have a repository
		//transaction, err := transactionService.CreateTransaction(context.Background(), transactionInput)
		//if err != nil {
		//	log.Println("❌ Error creating transaction:", err)
		//	return c.Send("❌ An error occurred while processing your transaction. Please try again later.")
		//}
		//balanceRepo := repository.NewBalanceRepoPG(db)
		//balanceService := services.NewBalanceService(balanceRepo)
		//_, err = balanceService.UpdateBalance(context.Background(), session.CampaignID, session.Amount)
		//if err != nil {
		//	log.Println("❌ Error updating balance:", err)
		//}

		// Update balance directly
		err = db.Model(&models.Balance{}).
			Where("campaign_id = ?", session.CampaignID).
			Update("amount", gorm.Expr("amount - ?", session.Amount)).Error
		if err != nil {
			log.Println("❌ Error updating balance directly:", err)
			// Not returning here, but you can return if it's critical
		}

		// Transaction created successfully
		log.Println("✅ Transaction created:", tx.TxReference)

		// Final success message
		delete(sessions, userID)
		return c.Send(fmt.Sprintf("🎉 Congratulations! Your coupon has been successfully redeemed.\nAmount paid: ₦%.2f\nThank you for participating!\n \n 👉 Follow @zidibot on Instagram, X & TikTok to join our next survey and win again!", session.Amount))
	}

	return nil
}

func saveCustomer(db *gorm.DB, customer *models.Customer) error {
	var existingCustomer models.Customer
	if err := db.Where("phone = ? AND campaign_id = ?", customer.Phone, customer.CampaignID).First(&existingCustomer).Error; err == nil {
		return fmt.Errorf("customer already exists for this campaign")
	}
	return db.Create(customer).Error
}

func saveResponses(db *gorm.DB, responses []models.Response) error {
	return db.Create(&responses).Error
}
