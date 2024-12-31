package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/middleware"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func User(app *fiber.App, db *gorm.DB) {
	userRepository := repositories.CommenceUserRepository(db)

	userService := services.CommenceUserService(userRepository)

	handler := handlers.UserHandler{IUserService: userService}

	general := app.Group("/v1/common")

	general.Get("/product", handler.GetProducts)
	general.Get("/product/:id", handler.GetProduct)

	user := app.Group("/v1/user")
	user.Use(middleware.ValidateJwt, middleware.UserRoleAuthentication)

	user.Post("/order", handler.CreateOrder)
	user.Get("/order", handler.GetOrders)
	user.Get("/order/:id", handler.GetOrder)
	user.Patch("", handler.UpdateUser)
	user.Patch("/order/:id", handler.UpdateOrder)
}
