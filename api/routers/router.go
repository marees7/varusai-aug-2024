package routers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RequiredRoute(app *fiber.App, db *gorm.DB) {
	AuthRoute(app, db)
	AdminRoute(app, db)
	UserRoute(app, db)
	MerchantRoute(app, db)
}
