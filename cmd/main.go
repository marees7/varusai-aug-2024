package main

import (
	"os"
	"shopping-site/api/routers"
	"shopping-site/internals"
	"shopping-site/pkg/loggers"

	"github.com/gofiber/fiber/v2"
)

func init() {
	internals.LoadEnvFile()
	loggers.ForLogs()
}

func main() {
	db := internals.InitiatePgConnection()
	internals.SchemaMigration(db)

	app := fiber.New()
	routers.RequiredRoute(app, db)

	err := app.Listen(os.Getenv("CLIENTPORT"))
	if err != nil {
		loggers.FatalLog.Fatal(err)
	}
}
