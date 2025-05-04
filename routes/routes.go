package routes

import (
	"github.com/Poojithkumar24/backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/backend")

	api.Post("/movie", handlers.AddMovie)
	api.Get("/movies", handlers.GetAllMovies)
}
