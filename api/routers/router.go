package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func RequiredRoute(app *fiber.App, db *gorm.DB) {
	// pass the db connection and fiber instace to all the routes
	Auth(app, db)
	Admin(app, db)
	User(app, db)
	Merchant(app, db)

	app.Get("/swagger/*", swagger.HandlerDefault)
}
