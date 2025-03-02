package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s Server) GetCampaignsCampaignIdQuestions(c *fiber.Ctx, campaignId openapi_types.UUID) error {
	if campaignId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid campaign ID",
		})
	}

	questions, err := s.questionService.GetQuestionsByCampaign(context.Background(), campaignId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Success",
		"Data":    questions,
	})
}

func (s Server) PostCampaignsCampaignIdQuestions(c *fiber.Ctx, campaignId openapi_types.UUID) error {
	if campaignId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid campaign ID",
		})
	}

	newQuestion := new(api.Question)
	if err := c.BodyParser(newQuestion); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   err.Error(),
		})
	}

	newQuestion.CampaignId = &campaignId

	createdQuestion, err := s.questionService.CreateQuestion(context.Background(), *newQuestion)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			ErrorCode: "500",
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdQuestion)
}

func (s Server) DeleteQuestionsQuestionId(c *fiber.Ctx, questionId openapi_types.UUID) error {
	if questionId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid question ID",
		})
	}

	err := s.questionService.DeleteQuestion(context.Background(), questionId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			ErrorCode: "404",
			Message:   "Question not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
