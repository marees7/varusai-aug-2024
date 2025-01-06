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
	// send db connection to repository
	userRepository := repositories.CommenceUserRepository(db)

	// send repo to service
	userService := services.CommenceUserService(userRepository)

	// Initialize the handler struct
	handler := handlers.UserHandler{IUserService: userService}

	// group the common endpoints
	general := app.Group("/v1/common")

	// common endpoints
	general.Get("/product", handler.GetProducts)
	general.Get("/product/:id", handler.GetProduct)

	// group the user endpoints
	user := app.Group("/v1/user")

	// added middleware for token validation and role authorization
	user.Use(middleware.ValidateJwt("user"), middleware.UserRoleAuthentication)

	// user endpoints
	user.Post("/order", handler.CreateOrder)
	user.Get("/order", handler.GetOrders)
	user.Get("/order/:id", handler.GetOrder)
	user.Patch("", handler.UpdateUser)
	user.Patch("/order/:id", handler.UpdateOrder)
}
