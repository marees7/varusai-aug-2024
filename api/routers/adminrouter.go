package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/middleware"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Admin(app *fiber.App, db *gorm.DB) {
	//send db connection to repository
	adminRepository := repositories.CommenceAdminRepository(db)

	// send repo to service
	adminService := services.CommenceAdminService(adminRepository)

	// Initialize the handler struct
	handler := handlers.AdminHandler{IAdminService: adminService}

	// group the admin endpoints
	admin := app.Group("/v1/admin")

	// added middleware for token validation and role authorization
	admin.Use(middleware.ValidateJwt("admin"), middleware.AdminRoleAuthentication)

	// admin endpoints
	admin.Post("/category", handler.CreateCategorey)
	admin.Post("/brand", handler.CreateBrand)
}
