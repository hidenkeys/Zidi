package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hidenkeys/zidibackend/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var jwtSecret = []byte("your_secret_key")

func (s Server) CreateOrganization(c *fiber.Ctx) error {
	organization := new(api.Organization)

	if err := c.BodyParser(organization); err != nil {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid request body",
		})
	}

	// Create the organization
	orgResponse, err := s.orgService.CreateOrganization(context.Background(), *organization)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	// Create an admin user for the organization
	defaultUser := api.User{
		Firstname:      "Admin",
		Lastname:       "User",
		Email:          organization.Email, // Assuming organization email is provided
		Password:       "ChangeMe123!",     // Default password (should be changed later)
		OrganizationId: orgResponse.Id,     // Assign newly created org ID
		Role:           "admin",
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   "Failed to encrypt password",
		})
	}
	defaultUser.Password = string(hashedPassword)

	userResponse, err := s.usrService.CreateUser(context.Background(), defaultUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	fmt.Println(userResponse.Id.String(), orgResponse.Id.String())

	// Generate a JWT Token
	token, err := generateJWTToken(userResponse.Id.String(), orgResponse.Id.String())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   "Failed to generate authentication token",
		})
	}

	// Return response with token
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"organization": orgResponse,
		"user":         userResponse,
		"token":        token,
	})
}

// Function to generate JWT Token
func generateJWTToken(userID, orgID string) (string, error) {
	claims := jwt.MapClaims{
		"userId":         userID,
		"organizationId": orgID,
		"exp":            time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s Server) GetOrganizations(c *fiber.Ctx) error {
	response, err := s.orgService.GetAllOrganizations(context.Background())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) DeleteOrganization(c *fiber.Ctx, organizationId openapi_types.UUID) error {
	err := s.orgService.DeleteOrganization(context.Background(), organizationId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusNoContent).JSON(nil)
}

func (s Server) GetOrganizationById(c *fiber.Ctx, organizationId openapi_types.UUID) error {
	response, err := s.orgService.GetOrganizationByID(context.Background(), organizationId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) UpdateOrganization(c *fiber.Ctx, organizationId openapi_types.UUID) error {
	organization := new(api.Organization)
	if err := c.BodyParser(organization); err != nil {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid request body",
		})
	}
	response, err := s.orgService.UpdateOrganization(context.Background(), organizationId, *organization)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (s Server) GetOrganizationByName(c *fiber.Ctx, params api.GetOrganizationByNameParams) error {
	response, err := s.orgService.GetOrganizationByName(context.Background(), params.Name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(response)
}
