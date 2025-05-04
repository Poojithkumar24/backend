package handlers

import (
	"context"

	db "github.com/Poojithkumar24/backend/db"

	"github.com/gofiber/fiber/v2"
)

func GetAllMovies(c *fiber.Ctx) error {
	var movies []db.Movie
	collection := db.MongoDB.Collection("movies")

	cursor, err := collection.Find(context.Background(), fiber.Map{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch movies"})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var movie db.Movie
		if err := cursor.Decode(&movie); err == nil {
			movies = append(movies, movie)
		}
	}

	return c.JSON(movies)
}
