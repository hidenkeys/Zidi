package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
)

type ResponseRepository interface {
	GetResponsesByQuestion(questionID uuid.UUID) ([]models.Response, error)
	CreateResponse(response *models.Response) (*models.Response, error)
}
