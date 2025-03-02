package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
