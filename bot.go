package main

// import (
// 	_ "encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/hidenkeys/zidibackend/config"
// 	"github.com/hidenkeys/zidibackend/models"
// 	"github.com/joho/godotenv"
// 	tele "gopkg.in/telebot.v3"
// 	"gorm.io/gorm"
// )

// type Session struct {
// 	CampaignID      uuid.UUID
// 	OrganizationID  uuid.UUID
// 	Amount          float64
// 	Customer        models.Customer
// 	Step            int
// 	Questions       []models.Question
// 	CurrentQuestion int
// 	Responses       []models.Response
// }

// var sessions = make(map[int64]*Session)

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("❌ Error loading .env file")
// 	}
// 	config.ConnectDatabase()

// 	bot, err := tele.NewBot(tele.Settings{
// 		Token:  os.Getenv("TELEGRAM_API_KEY"),
// 		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
// 	})
// 	if err != nil {
// 		log.Fatal("❌ Error initializing bot: ", err)
// 	}

// 	log.Println("🚀 Bot is running...")

// 	bot.Handle("/start", func(c tele.Context) error {
// 		args := c.Args()
// 		if len(args) == 0 {
// 			return c.Send("❌ Please provide a valid campaign ID.")
// 		}

// 		campaignID, err := uuid.Parse(args[0])
// 		if err != nil {
// 			return c.Send("❌ Invalid campaign ID format.")
// 		}

// 		var campaign models.Campaign
// 		if err := config.DB.Where("id = ? AND status = ?", campaignID, "active").First(&campaign).Error; err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				return c.Send("❌ Campaign not found or not active.")
// 			}
// 			return c.Send("❌ An error occurred while fetching the campaign.")
// 		}

// 		var questions []models.Question
// 		if err := config.DB.Where("campaign_id = ?", campaignID).Find(&questions).Error; err != nil {
// 			return c.Send("❌ Failed to fetch questions.")
// 		}

// 		user := c.Sender()

// 		// Check if the customer already exists for the same campaign
// 		var existingCustomer models.Customer
// 		if err := config.DB.Where("phone = ? AND campaign_id = ?", user.Recipient(), campaign.ID).First(&existingCustomer).Error; err == nil {
// 			return c.Send("✅ You have already registered for this campaign. Stay tuned for the next one!")
// 		}

// 		sessions[user.ID] = &Session{
// 			CampaignID:     campaign.ID,
// 			OrganizationID: campaign.OrganizationID,
// 			Amount:         campaign.Amount,
// 			Customer: models.Customer{
// 				Status:     "inactive",
// 				CampaignID: campaign.ID,
// 			},
// 			Step:            1,
// 			Questions:       questions,
// 			CurrentQuestion: 0,
// 			Responses:       []models.Response{},
// 		}

// 		return c.Send(fmt.Sprintf("👋 Welcome!\n%s\n\nLet's get started by gathering your details.\nWhat's your first name?", campaign.WelcomeMessage))
// 	})

// 	bot.Handle(tele.OnText, func(c tele.Context) error {
// 		userID := c.Sender().ID
// 		session, exists := sessions[userID]
// 		if !exists {
// 			return c.Send("❌ Please start with /start and a valid campaign ID.")
// 		}

// 		switch session.Step {
// 		case 1:
// 			session.Customer.FirstName = c.Text()
// 			session.Step++
// 			return c.Send("📛 What's your last name?")
// 		case 2:
// 			session.Customer.LastName = c.Text()
// 			session.Step++
// 			return c.Send("📧 What’s your email address?")
// 		case 3:
// 			session.Customer.Email = c.Text()
// 			session.Step++
// 			return c.Send("📞 Please provide your phone number:")
// 		case 4:
// 			session.Customer.Phone = c.Text()
// 			session.Step++
// 			return c.Send("📶 Which network provider do you use?")
// 		case 5:
// 			session.Customer.Network = c.Text()
// 			session.Step++
// 			return c.Send("💬 Any feedback you'd like to share?")
// 		case 6:
// 			session.Customer.Feedback = c.Text()
// 			session.Customer.OrganizationID = session.OrganizationID
// 			session.Customer.Amount = session.Amount

// 			if err := saveCustomer(&session.Customer); err != nil {
// 				log.Println("❌ Error saving customer:", err)
// 				return c.Send("❌ An error occurred while saving your details. Please try again.")
// 			}

// 			if len(session.Questions) > 0 {
// 				session.Step = 7
// 				return c.Send(fmt.Sprintf("📋 Now, let's answer %d additional questions.\n%s", len(session.Questions), session.Questions[0].Text))
// 			}

// 			delete(sessions, userID)
// 			return c.Send("🎉 Thank you! Your details have been successfully saved.")

// 		case 7:
// 			question := session.Questions[session.CurrentQuestion]
// 			session.Responses = append(session.Responses, models.Response{
// 				CustomerID: session.Customer.ID,
// 				QuestionID: question.ID,
// 				Answer:     c.Text(),
// 			})

// 			if session.CurrentQuestion+1 < len(session.Questions) {
// 				session.CurrentQuestion++
// 				return c.Send(session.Questions[session.CurrentQuestion].Text)
// 			}

// 			if err := saveResponses(session.Responses); err != nil {
// 				log.Println("❌ Error saving responses:", err)
// 				return c.Send("❌ An error occurred while saving your responses.")
// 			}

// 			delete(sessions, userID)
// 			return c.Send("🎉 Thank you! Your details and responses have been successfully saved.")
// 		}
// 		return nil
// 	})

// 	bot.Start()
// }

// func saveCustomer(customer *models.Customer) error {
// 	var existingCustomer models.Customer
// 	if err := config.DB.Where("phone = ? AND campaign_id = ?", customer.Phone, customer.CampaignID).First(&existingCustomer).Error; err == nil {
// 		return fmt.Errorf("customer already exists for this campaign")
// 	}
// 	return config.DB.Create(customer).Error
// }

func saveResponses(responses []models.Response) error {
	return config.DB.Create(&responses).Error
}
