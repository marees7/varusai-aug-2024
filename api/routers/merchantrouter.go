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
	// send db connection to repository
	merchantRepository := repositories.CommenceMerchantRepository(db)

	// send repo to service
	merchantService := services.CommenceMerchantService(merchantRepository)

	// Initialize the handler struct
	handler := handlers.MerchantHandler{IMerchantService: merchantService}

	// group the merchant endpoints
	merchant := app.Group("/v1/merchant")

	// added middleware for token validation and role authorization
	merchant.Use(middleware.ValidateJwt("merchant"), middleware.MerchantRoleAuthentication)

	// merchant endpoints
	merchant.Post("/product", handler.CreateProduct)
	merchant.Get("/category", handler.GetCategories)
	merchant.Get("/brand", handler.GetBrands)
	merchant.Get("product", handler.GetProducts)
	merchant.Get("/product/:id", handler.GetProduct)
	merchant.Get("/order", handler.GetOrders)
	merchant.Get("/order/:id", handler.GetOrder)
	merchant.Patch("/product", handler.UpdateProduct)
	merchant.Patch("", handler.UpdateMerchant)
	merchant.Patch("/order/:id", handler.UpdateOrderStatus)
	merchant.Delete("/product/:id", handler.DeleteProduct)
}
