package main

import (
	"log"

	"github.com/Poojithkumar24/backend/routes"

	db "github.com/Poojithkumar24/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,OPTIONS",
	}))
	db.ConnectDB()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":4000"))
}
