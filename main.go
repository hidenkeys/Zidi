package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/config"
	"github.com/hidenkeys/zidibackend/handlers"
	"github.com/hidenkeys/zidibackend/repository"
	"github.com/hidenkeys/zidibackend/services"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	config.ConnectDatabase()
	config.MigrateDatabase()

	db := config.DB

	orgRepo := repository.NewOrganizationRepoPG(db)
	orgService := services.NewOrganisationService(orgRepo)

	userRepo := repository.NewUserRepoPG(db)
	userService := services.NewUserService(userRepo)

	customerRepo := repository.NewCustomerRepoPG(db)
	customerService := services.NewCustomerService(customerRepo)

	campaignRepo := repository.NewCampaignRepoPG(db)
	campaignService := services.NewCampaignService(campaignRepo)

	questionRepo := repository.NewQuestionRepoPG(db)
	questionService := services.NewQuestionService(questionRepo)

	responseRepo := repository.NewResponseRepoPG(db)
	responseService := services.NewResponseService(responseRepo)

	server := handlers.NewServer(orgService, userService, campaignService, customerService, questionService, responseService)

	app := fiber.New()

	api.RegisterHandlers(app, server)

	// And we serve HTTP until the world ends.
	log.Fatal(app.Listen("0.0.0.0:8080"))

}
