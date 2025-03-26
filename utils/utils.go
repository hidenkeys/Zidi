package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type UserClaims struct {
	ID             string `json:"user_id"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
}

// RevokeToken stores a token in the database
func RevokeToken(db *gorm.DB, token string, expiry time.Time) error {
	t := models.Token{
		Token:     token,
		ExpiresAt: expiry,
	}
	return db.Create(&t).Error
}

func IsTokenRevoked(db *gorm.DB, token string) (bool, error) {
	var count int64
	err := db.Model(&models.Token{}).Where("token = ?", token).Count(&count).Error
	return count > 0, err
}

var jwtSecret = []byte("your-secret-key")

func GenerateJWTToken(userID, orgID, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId":         userID,
		"organizationId": orgID,
		"role":           role,                                  // Include role in JWT
		"exp":            time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateTokens(charType string, length int, count int) []string {
	charsets := map[string]string{
		"alphanumerical": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"numeric":        "0123456789",
	}

	charset, ok := charsets[charType]
	if !ok {
		charset = charsets["alphanumerical"] // Default to alphanumerical
	}

	tokens := make(map[string]struct{})
	var tokenList []string

	for len(tokenList) < count {
		token := make([]byte, length)
		for i := range token {
			token[i] = charset[rand.Intn(len(charset))]
		}

		tokenStr := string(token)
		if _, exists := tokens[tokenStr]; !exists {
			tokens[tokenStr] = struct{}{}
			tokenList = append(tokenList, tokenStr)
		}
	}

	return tokenList
}

// UserRoles defines user roles
const (
	RoleZidi  = "zidi"
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// CheckUserRole checks if the user has a required role
func CheckUserRole(c *fiber.Ctx, allowedRoles ...string) bool {
	user, ok := c.Locals("user").(UserClaims)
	if !ok {
		return false
	}

	// Convert role to lowercase and compare
	userRole := strings.ToLower(user.Role)
	for _, role := range allowedRoles {
		if userRole == strings.ToLower(role) {
			return true
		}
	}
	return false
}

// function to send email
func SendEmail(to, subject, body string) error {
	from := "teniolasobande04@gmail.com"
	password := "vndt vleo ccfc tcqt"

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}
	return nil
}

// PaystackResponse holds the response from Paystack API
type PaystackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

// CreatePaystackPaymentLink initializes a payment with metadata
func CreatePaystackPaymentLink(email string, amount int, campaignID, organizationID string) (string, error) {
	paystackURL := "https://api.paystack.co/transaction/initialize"
	apiKey := os.Getenv("PAYSTACK_SK") // Store in .env

	// Convert amount to kobo (Paystack requires amount in kobo)
	amountKobo := amount * 100

	// Create the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"email":  email,
		"amount": amountKobo,
		"metadata": map[string]interface{}{
			"campaign_id":     campaignID,
			"organization_id": organizationID,
		},
	})
	if err != nil {
		return "", err
	}

	// Make the HTTP request
	req, err := http.NewRequest("POST", paystackURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Parse the response
	var paystackResp PaystackResponse
	if err := json.NewDecoder(resp.Body).Decode(&paystackResp); err != nil {
		return "", err
	}

	if !paystackResp.Status {
		return "", fmt.Errorf("failed to initialize payment: %s", paystackResp.Message)
	}

	return paystackResp.Data.AuthorizationURL, nil
}

// FlutterwaveResponse defines the response structure
type FlutterwaveResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

func CreateFlutterwavePaymentLink(email string, amount int, campaignID, organizationID string) (string, error) {
	flutterwaveURL := "https://api.flutterwave.com/v3/payments"
	apiKey := os.Getenv("FLW_SECRET_KEY")

	// Generate a unique transaction reference
	txRef := fmt.Sprintf("%s-%s", campaignID, uuid.New().String())

	requestBody, err := json.Marshal(map[string]interface{}{
		"tx_ref":       txRef,
		"amount":       amount,
		"currency":     "NGN",
		"redirect_url": "https://yourwebsite.com/payment-success",
		"customer": map[string]string{
			"email": email,
		},
		"metadata": map[string]string{
			"campaign_id":     campaignID,
			"organization_id": organizationID,
		},
		"customizations": map[string]string{
			"title":       "Campaign Payment",
			"description": "Payment for campaign",
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", flutterwaveURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var flutterwaveResp FlutterwaveResponse
	if err := json.NewDecoder(resp.Body).Decode(&flutterwaveResp); err != nil {
		return "", err
	}

	if flutterwaveResp.Status != "success" {
		return "", fmt.Errorf("failed to initialize payment: %s", flutterwaveResp.Message)
	}

	return flutterwaveResp.Data.Link, nil
}

func VerifyFlutterwaveSignature(body []byte, signature, secret string) bool {
	fmt.Println("00")
	if secret == "" || signature == "" {
		return false
	}

	fmt.Println("01")

	// Compute HMAC SHA-256 hash of the request body
	hash := hmac.New(sha256.New, []byte(secret))
	fmt.Println("03")
	hash.Write(body)
	fmt.Println("03")
	expectedSignature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	fmt.Println("signature", signature)
	fmt.Println("expected sigaturew", expectedSignature)
	fmt.Println("04")

	return signature == expectedSignature
}

// verifyFlutterwaveTransaction checks transaction status from Flutterwave API
func VerifyFlutterwaveTransaction(transactionID int) (bool, error) {
	apiURL := fmt.Sprintf("https://api.flutterwave.com/v3/transactions/%d/verify", transactionID)
	apiKey := os.Getenv("FLW_SECRET_KEY") // Get secret key from .env

	// Make request to Flutterwave API
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Parse JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	// Check if the transaction was successful
	if status, ok := response["status"].(string); ok && status == "success" {
		data, exists := response["data"].(map[string]interface{})
		if exists {
			if txStatus, ok := data["status"].(string); ok && txStatus == "successful" {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("transaction not successful")
}
