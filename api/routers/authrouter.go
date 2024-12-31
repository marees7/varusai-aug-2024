package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Auth(app *fiber.App, db *gorm.DB) {
	authRepository := repositories.CommenceAuthRepository(db)

	authService := services.CommenceAuthService(authRepository)

	handler := handlers.AuthHandler{IAuthService: authService}

	app.Post("/signup", handler.Signup)
	app.Post("/login", handler.Login)
}
