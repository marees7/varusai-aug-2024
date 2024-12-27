package routers

import (
	"shopping-site/api/handlers"
	"shopping-site/api/middleware"
	"shopping-site/api/repositories"
	"shopping-site/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AdminRoute(app *fiber.App, db *gorm.DB) {
	adminRepository := repositories.CommenceAdminRepository(db)

	adminService := services.CommenceAdminService(adminRepository)

	handler := handlers.AdminHandler{IAdminService: adminService}

	user := app.Group("/v1/role/admin")
	user.Use(middleware.ValidateJwt, middleware.AdminRoleAuthentication)

	user.Post("/category", handler.AddCategoreyHandler)
	user.Post("/brand", handler.AddBrandHandler)

}
