package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/models"
	"gorm.io/gorm"
)

type questionPG struct {
	db *gorm.DB
}

func NewQuestionRepoPG(db *gorm.DB) QuestionRepository {
	return &questionPG{db: db}
}

func (r *questionPG) GetQuestionsByCampaign(campaignID uuid.UUID) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Where("campaign_id = ?", campaignID).Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *questionPG) CreateQuestion(question *models.Question) (*models.Question, error) {
	fmt.Println("this is question inn my question pg", question)
	err := r.db.Create(question).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("this is question after create in my question pg", question)
	return question, nil
}

func (r *questionPG) DeleteQuestion(questionID uuid.UUID) error {
	err := r.db.Where("id = ?", questionID).Delete(&models.Question{}).Error
	if err != nil {
		return err
	}
	return nil
}
