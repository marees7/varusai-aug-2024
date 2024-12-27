package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/middleware"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MerchantRoute(app *fiber.App, db *gorm.DB) {
	merchantRepository := repositories.CommenceMerchantRepository(db)

	merchantService := services.CommenceMerchantService(merchantRepository)

	handler := handlers.MerchantHandler{IMerchantService: merchantService}

	merchant := app.Group("/v1/role/merchant")
	merchant.Use(middleware.ValidateJwt, middleware.MerchantRoleAuthentication)

	merchant.Post("/product", handler.AddProductHandler)
	merchant.Get("product", handler.GetProductsHandler)
	merchant.Get("/order", handler.GetOrdersHandler)
	merchant.Get("/product/:id", handler.GetProductHandler)
	merchant.Patch("/product", handler.UpdateProductHandler)
	merchant.Patch("", handler.UpdateMerchantHandler)
	merchant.Patch("/order:id", handler.UpdateOrderStatusHandler)
	merchant.Delete("/product/:id", handler.RemoveProductHandler)
}
