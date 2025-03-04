package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

var jwtSecret = []byte("your-secret-key")

func GenerateJWT(userID uuid.UUID, role string, orgID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id":         userID.String(),
		"role":            role,
		"organization_id": orgID.String(),
		"exp":             time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":             time.Now().Unix(),
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
