package handlers

import (
	"github.com/hidenkeys/zidibackend/services"
)

type Server struct {
	orgService      *services.OrganizationService
	usrService      *services.UserService
	campaignService *services.CampaignService
	customerService *services.CustomerService
	questionService *services.QuestionService
	responseService *services.ResponseService
}

func NewServer(orgService *services.OrganizationService, usrService *services.UserService, campaignService *services.CampaignService, customerService *services.CustomerService, questionService *services.QuestionService, responseService *services.ResponseService) *Server {
	return &Server{
		orgService:      orgService,
		usrService:      usrService,
		campaignService: campaignService,
		customerService: customerService,
		questionService: questionService,
		responseService: responseService,
	}
}
