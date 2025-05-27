// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"

// 	"github.com/google/uuid"
// 	"github.com/joho/godotenv"
// 	"github.com/xuri/excelize/v2"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"

// 	"github.com/hidenkeys/zidibackend/models"
// 	"github.com/hidenkeys/zidibackend/utils"
// )

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	dsn := os.Getenv("DATABASE_URL")
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	f, err := excelize.OpenFile("Zidibackend_public_customers.xlsx")
// 	if err != nil {
// 		log.Fatal("Error opening Excel file:", err)
// 	}

// 	rows, err := f.GetRows("Result 1") // Adjust sheet name if needed
// 	if err != nil {
// 		log.Fatal("Error reading sheet:", err)
// 	}

// 	campaignID := uuid.MustParse("082432db-b05b-48e8-9201-b06f2b9254dd")
// 	orgID := uuid.MustParse("41c9d707-58d4-4a6e-b18d-29ad7865ee8e")
// 	amount := 100.0

// 	for idx, row := range rows {
// 		if idx == 0 || len(row) < 4 {
// 			continue // skip header or incomplete row
// 		}

// 		name := row[4] + " " + row[5]                         // Combine first and last name
// 		network := strings.ToLower(strings.TrimSpace(row[9])) // Correct column for network
// 		phone := sanitizePhone(row[6])                        // Correct column for phone

// 		fmt.Printf("\nProcessing: %s - %s\n", name, phone)

// 		var customer models.Customer
// 		err := db.Where("phone = ? AND campaign_id = ?", phone, campaignID).First(&customer).Error
// 		if err != nil {
// 			log.Printf("❌ Customer not found for %s: %v\n", phone, err)
// 			continue
// 		}

// 		var status, reference string
// 		var commission float64

// 		airtimeRes, err := utils.SendAirtime(fmt.Sprintf("%.0f", amount), network, phone)
// 		if err != nil {
// 			log.Printf("⚠️ Airtime failed for %s: %v\n", phone, err)
// 			status = "failed"
// 		} else {
// 			status = airtimeRes.Status
// 			reference = airtimeRes.RequestID
// 			commission, _ = strconv.ParseFloat(airtimeRes.Commission, 64)

// 			// Only update balance and status if successful
// 			if status == "successful" {
// 				err = db.Model(&models.Balance{}).
// 					Where("campaign_id = ?", campaignID).
// 					Update("amount", gorm.Expr("amount - ?", amount)).Error
// 				if err != nil {
// 					log.Printf("❌ Balance update failed for %s: %v\n", phone, err)
// 				}

// 				_ = db.Model(&models.Customer{}).Where("id = ?", customer.ID).Update("status", "active")
// 			}
// 		}

// 		tx := models.Transaction{
// 			OrganizationID: orgID,
// 			CampaignID:     campaignID,
// 			CustomerID:     customer.ID,
// 			Network:        network,
// 			PhoneNumber:    phone,
// 			TxReference:    reference,
// 			Status:         status,
// 			Amount:         amount,
// 			Type:           "airtime",
// 			Commisson:      commission,
// 		}
// 		db.Create(&tx)
// 	}
// }

// func sanitizePhone(input string) string {
// 	trimmed := strings.TrimSpace(input)
// 	trimmed = strings.ReplaceAll(trimmed, " ", "")
// 	trimmed = strings.ReplaceAll(trimmed, "-", "")
// 	return trimmed
// }
