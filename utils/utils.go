package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
	defer resp.Body.Close()

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
