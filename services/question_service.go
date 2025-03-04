package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"

	"github.com/google/uuid"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/models"
	"github.com/hidenkeys/zidibackend/repository"
)

type QuestionService struct {
	questionRepo repository.QuestionRepository
}

func NewQuestionService(questionRepo repository.QuestionRepository) *QuestionService {
	return &QuestionService{questionRepo: questionRepo}
}

func (s *QuestionService) CreateQuestion(ctx context.Context, req api.Question) (*api.Question, error) {
	var optionsJSON datatypes.JSON
	if req.Options != nil {
		jsonData, err := json.Marshal(req.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal options: %w", err)
		}
		optionsJSON = datatypes.JSON(jsonData)
	}
	newQuestion := &models.Question{
		ID:         req.Id,
		CampaignID: req.CampaignId,
		Text:       req.Text,
		Type:       string(req.Type),
		Options:    optionsJSON,
	}

	question, err := s.questionRepo.CreateQuestion(newQuestion)
	if err != nil {
		return nil, err
	}

	fmt.Println("this is my question", question)
	return mapToAPIQuestion(question), nil
}

func (s *QuestionService) GetQuestionsByCampaign(ctx context.Context, campaignID uuid.UUID) ([]api.Question, error) {
	questions, err := s.questionRepo.GetQuestionsByCampaign(campaignID)
	if err != nil {
		return nil, err
	}

	var finalQuestions []api.Question
	for _, question := range questions {
		finalQuestions = append(finalQuestions, *mapToAPIQuestion(&question))
	}

	return finalQuestions, nil
}

//func (s *QuestionService) GetQuestionByID(ctx context.Context, id uuid.UUID) (*api.Question, error) {
//	question, err := s.questionRepo.(id)
//	if err != nil {
//		return nil, err
//	}
//
//	return mapToAPIQuestion(&question), nil
//}

func (s *QuestionService) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return s.questionRepo.DeleteQuestion(id)
}

// Helper function to convert models.Question to api.Question
func mapToAPIQuestion(question *models.Question) *api.Question {
	var options []string
	if question.Options != nil {
		_ = json.Unmarshal(question.Options, &options) // Convert JSONB to []string
	}

	return &api.Question{
		Id:         question.ID,
		CampaignId: question.CampaignID,
		Text:       question.Text,
		Type:       api.QuestionType(question.Type),
		Options:    &options,
	}
}
