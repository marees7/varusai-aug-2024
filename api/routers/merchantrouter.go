package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/middleware"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Merchant(app *fiber.App, db *gorm.DB) {
	merchantRepository := repositories.CommenceMerchantRepository(db)

	merchantService := services.CommenceMerchantService(merchantRepository)

	handler := handlers.MerchantHandler{IMerchantService: merchantService}

	merchant := app.Group("/v1/merchant")
	merchant.Use(middleware.ValidateJwt, middleware.MerchantRoleAuthentication)

	merchant.Post("/product", handler.CreateProduct)
	merchant.Get("product", handler.GetProducts)
	merchant.Get("/product/:id", handler.GetProduct)
	merchant.Get("/order", handler.GetOrders)
	merchant.Get("/order/:id", handler.GetOrder)
	merchant.Patch("/product", handler.UpdateProduct)
	merchant.Patch("", handler.UpdateMerchant)
	merchant.Patch("/order/:id", handler.UpdateOrderStatus)
	merchant.Delete("/product/:id", handler.DeleteProduct)
}
