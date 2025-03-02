package repository

import (
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
)

type responsePG struct {
	db *gorm.DB
}

func NewResponseRepoPG(db *gorm.DB) ResponseRepository {
	return &responsePG{db: db}
}

func (r *responsePG) GetResponsesByQuestion(questionID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := r.db.Where("question_id = ?", questionID).Find(&responses).Error
	return responses, err
}

func (r *responsePG) CreateResponse(response *models.Response) (*models.Response, error) {
	err := r.db.Create(response).Error
	if err != nil {
		return nil, err
	}
	return response, nil
}
