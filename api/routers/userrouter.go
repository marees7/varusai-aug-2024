package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/middleware"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoute(app *fiber.App, db *gorm.DB) {
	userRepository := repositories.CommenceUserRepository(db)

	userService := services.CommenceUserService(userRepository)

	handler := handlers.UserHandler{IUserService: userService}

	user := app.Group("/v1/role/user")
	user.Use(middleware.ValidateJwt)

	user.Post("/order", handler.PlaceOrderHandler)
	user.Get("/order", handler.GetOrdersHandler)
	user.Get("/product/filter", handler.FilterProductsHandler)
	user.Get("product", handler.GetProductsHandler)
	user.Get("/product/:id", handler.GetProductHandler)
	user.Patch("", handler.UpdateUserHandler)
	user.Patch("/order/:id", handler.CancelOrderHandler)
}
