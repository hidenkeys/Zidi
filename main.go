package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hidenkeys/zidibackend/api"
	"github.com/hidenkeys/zidibackend/config"
	"github.com/hidenkeys/zidibackend/handlers"
	"github.com/hidenkeys/zidibackend/middleware"
	"github.com/hidenkeys/zidibackend/repository"
	"github.com/hidenkeys/zidibackend/services"
	"github.com/hidenkeys/zidibackend/telegrambot"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	var jwtSecret = []byte("your-secret-key")
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://www.zidi-admin.vercel.app, https://zidi-admin.vercel.app, https://zidi-frontend.vercel.app, https://zidi-frontend.vercel.app/, https://216.198.79.65:3000, https://64.29.17.65:3000 ", // Allow your frontend origin
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",                                                                                                                                                                                     // Allow specific HTTP methods
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",                                                                                                                                                                         // Allow custom headers
	}))
	userAuth := middleware.AuthMiddleware(string(jwtSecret), "user", "admin", "zidi")
	app.Post("/api/v1/auth/login", server.LoginUser)
	app.Use(userAuth)

	//adminAuth := middleware.AuthMiddleware(string(jwtSecret), "admin")
	//zidiAuth := middleware.AuthMiddleware(string(jwtSecret), "zidi")
	//zidiAndAdminAuth := middleware.AuthMiddleware(string(jwtSecret), "zidi","admin")
	//adminAndUserAuth := middleware.AuthMiddleware(string(jwtSecret), )
	go telegrambot.StartBot(db)

	api.RegisterHandlersWithOptions(app, server, api.FiberServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: []api.MiddlewareFunc{
			//userAuth,
		},
	})

	server.SeedDefaultOrganization()
	// And we serve HTTP until the world ends.
	log.Fatal(app.Listen("0.0.0.0:8080"))

}
