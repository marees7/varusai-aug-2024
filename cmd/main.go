package main

import (
	"os"
	"shopping-site/api/routers"
	"shopping-site/internals"
	"shopping-site/pkg/loggers"

	_ "shopping-site/docs"

	"github.com/gofiber/fiber/v2"
)

func init() {
	// load env file
	internals.LoadEnvFile()
	// initialize loggers
	loggers.ForLogs()
}

//	@title			Shopping-site API
//	@version		1.0
//	@description	It's an basic ecommerce site where you can view products and order any product. Then also you can be a seller by listing your products
//
// @contact.name	varusai
// @contact.email	varusi0605@gmail.com
// @license.url	https://github.com/marees7/varusai-aug-2024
// @host			localhost:8080
// @BasePath		/v1
func main() {
	// initiate db and connection
	db := internals.InitiatePgConnection()
	// migrate the models for tables
	internals.SchemaMigration(db)

	// initiate new fiber router and get instance
	app := fiber.New()
	// pass the fiber instace and db connection to the routes
	routers.RequiredRoute(app, db)

	// start the server and listen for request
	err := app.Listen(os.Getenv("PORT"))
	if err != nil {
		loggers.FatalLog.Fatal(err)
	}
}
