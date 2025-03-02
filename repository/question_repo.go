package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
)

type QuestionRepository interface {
	GetQuestionsByCampaign(campaignID uuid.UUID) ([]models.Question, error)
	CreateQuestion(question *models.Question) (*models.Question, error)
	DeleteQuestion(questionID uuid.UUID) error
}
