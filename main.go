package main

import (
	"log"

	"github.com/Poojithkumar24/backend/routes"

	db "github.com/Poojithkumar24/backend/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	db.ConnectDB()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
