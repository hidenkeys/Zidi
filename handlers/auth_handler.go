package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s Server) LoginUser(c *fiber.Ctx) error {
	var req api.LoginUserJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Email or Password is required",
		})
	}

	user, err := s.usrService.GetUserByEmail(context.Background(), string(req.Email))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(api.Error{
			ErrorCode: "401",
			Message:   err.Error(),
		})
	}

	token, err := utils.GenerateJWT(user.Id, user.Role, user.OrganizationId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	user.Password = ""

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"Message": "Success",
		"token":   token,
		"user":    user,
	})
}
