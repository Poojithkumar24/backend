package handlers

import (
	"context"

	db "github.com/Poojithkumar24/backend/db"

	"github.com/gofiber/fiber/v2"
)

func AddMovie(c *fiber.Ctx) error {
	var movie db.Movie

	// Parse request body into the movie struct
	if err := c.BodyParser(&movie); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// Insert into MongoDB
	collection := db.MongoDB.Collection("movies")
	result, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert movie",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "movie added successfully",
		"id":      result.InsertedID,
	})
}

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
