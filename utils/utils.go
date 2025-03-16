package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"math/rand"
	"net/smtp"
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
