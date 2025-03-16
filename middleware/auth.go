package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// UserClaims represents extracted claims from the JWT
type UserClaims struct {
	ID             string `json:"user_id"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
}

// Claims structure for JWT parsing
type Claims struct {
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
	jwt.RegisteredClaims
}

// AllowedRoles defines the roles permitted to access routes
var AllowedRoles = map[string]bool{
	"zidi":  true,
	"admin": true,
	"user":  true,
}

// AuthMiddleware validates the JWT and enforces role-based access
func AuthMiddleware(secretKey string, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Extract claims properly
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Check if the token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
		}

		// Validate role from the allowed roles
		if !AllowedRoles[claims.Role] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized role"})
		}

		// Check if the user's role is in the required roles for this route
		if len(allowedRoles) > 0 {
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
			}
		}

		// Attach claims to context for later use
		c.Locals("user", UserClaims{
			ID:             claims.UserID,
			Role:           claims.Role,
			OrganizationID: claims.OrganizationID,
		})

		return c.Next()
	}
}
