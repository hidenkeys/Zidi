package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/config"
	"github.com/hidenkeys/zidibackend/models"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

// Step tracking for user sessions
type Session struct {
	CampaignID     uuid.UUID
	OrganizationID uuid.UUID
	Amount         float64
	Customer       models.Customer
	Step           int
}

var sessions = make(map[int64]*Session)

func main() {
	// Load environment variables and connect to the database
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDatabase()

	// Initialize the bot
	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_API_KEY"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Handle /start command with campaign ID
	bot.Handle("/start", func(c tele.Context) error {
		args := c.Args()
		if len(args) == 0 {
			return c.Send("❌ Please provide a valid campaign ID.")
		}

		// Parse the campaign ID
		campaignID, err := uuid.Parse(args[0])
		if err != nil {
			return c.Send("❌ Invalid campaign ID format.")
		}

		// Find the campaign by ID
		var campaign models.Campaign
		if err := config.DB.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Send("❌ Campaign not found.")
			}
			log.Println("Error fetching campaign:", err)
			return c.Send("❌ Something went wrong.")
		}

		// Start a new session for the user
		user := c.Sender()
		sessions[user.ID] = &Session{
			CampaignID:     campaign.ID,
			OrganizationID: campaign.OrganizationID,
			Amount:         campaign.Amount,
			Customer: models.Customer{
				Status:     "inactive",
				CampaignID: campaign.ID,
			},
			Step: 1,
		}

		// Send the welcome message
		return c.Send(fmt.Sprintf("👋 Welcome!\n%s\n\nLet's get started by gathering your details.\nwhat's your first name?", campaign.WelcomeMessage))
	})

	// Handle user responses and collect customer information
	bot.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Sender().ID
		session, exists := sessions[userID]
		if !exists {
			return c.Send("❌ Please start the process using /start with a campaign ID.")
		}

		// Handle each question in sequence
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
			return c.Send("📞 Please provide your phone number:")
		case 4:
			session.Customer.Phone = c.Text()
			session.Step++
			return c.Send("📶 Which network provider do you use?")
		case 5:
			session.Customer.Network = c.Text()
			session.Step++
			return c.Send("💬 Any feedback you'd like to share?")
		case 6:
			session.Customer.Feedback = c.Text()
			session.Step++
			fmt.Println("this is the customer", session.Customer)
			return c.Send("✅ All information collected. Saving your details...")

		case 7:
			// Set additional fields from the session
			session.Customer.OrganizationID = session.OrganizationID
			session.Customer.Amount = session.Amount

			// Save customer information
			if err := saveCustomer(session.Customer); err != nil {
				log.Println("Error saving customer:", err)
				return c.Send("❌ An error occurred while saving your details. Please try again.")
			}

			delete(sessions, userID) // Clear session after saving
			return c.Send("🎉 Thank you! Your details have been successfully saved.")
		}

		return nil
	})

	log.Println("🚀 Bot is running...")
	bot.Start()
}

// Save the customer to the database
func saveCustomer(customer models.Customer) error {
	var existingCustomer models.Customer
	result := config.DB.Where("phone = ?", customer.Phone).First(&existingCustomer)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new customer if not found
		return config.DB.Create(&customer).Error
	}

	// Update customer if exists
	existingCustomer.FirstName = customer.FirstName
	existingCustomer.LastName = customer.LastName
	existingCustomer.Email = customer.Email
	existingCustomer.Network = customer.Network
	existingCustomer.Feedback = customer.Feedback
	existingCustomer.Amount = customer.Amount
	existingCustomer.CampaignID = customer.CampaignID
	existingCustomer.OrganizationID = customer.OrganizationID
	existingCustomer.Status = "inactive"

	return config.DB.Save(&existingCustomer).Error
}
