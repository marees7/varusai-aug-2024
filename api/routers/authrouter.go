package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	// send the db connection to repository
	authRepository := repositories.CommenceAuthRepository(db)

	// send the repo  to service
	authService := services.CommenceAuthService(authRepository)

	// Initialize the handler struct
	handler := handlers.AuthHandler{IAuthService: authService}

	//auth endpoints
	app.Post("v1/signup", handler.Signup)
	app.Post("v1/login", handler.Login)
}
