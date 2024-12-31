package routers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RequiredRoute(app *fiber.App, db *gorm.DB) {
	Auth(app, db)
	Admin(app, db)
	User(app, db)
	Merchant(app, db)
}
