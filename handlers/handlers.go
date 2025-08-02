package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	db "github.com/Poojithkumar24/backend/db"
	"github.com/Poojithkumar24/backend/tmdb"

	"github.com/gofiber/fiber/v2"
)

func AddMovie(c *fiber.Ctx) error {

	log.Println("add movie called ")
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

// func GetAllMovies(c *fiber.Ctx) error {
// 	var movies []db.Movie
// 	collection := db.MongoDB.Collection("movies")

// 	cursor, err := collection.Find(context.Background(), fiber.Map{})
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch movies"})
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var movie db.Movie
// 		if err := cursor.Decode(&movie); err == nil {
// 			movies = append(movies, movie)
// 		}
// 	}

// 	return c.JSON(movies)
// }

func GetAllMovies(c *fiber.Ctx) error {
	apiKey := "706f38368bcaff89d9e429410278545f"
	url := "https://api.themoviedb.org/3/movie/upcoming?api_key=" + apiKey + "&page=1"

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch data from TMDb"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to read response body"})
	}

	var apiResp tmdb.TMDbResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to parse TMDb response"})
	}

	var movies []db.Movie
	for _, tmdbMovie := range apiResp.Results {
		movie := db.Movie{
			MovieID:     "", // TMDb doesn't return a string ID; you can optionally use strconv.Itoa(tmdbMovie.ID)
			Title:       tmdbMovie.Title,
			Genre:       "Unknown", // TMDb `genre_ids` is an array of IDs — need extra API call to resolve
			Duration:    0,         // Not included in `/movie/popular` — requires another API call
			Description: tmdbMovie.Overview,
			Screen:      "Main Screen",                       // Default/static value
			Shows:       []string{"10:00", "13:00", "18:00"}, // Default showtimes
			ImageUrl:    "https://image.tmdb.org/t/p/w780" + tmdbMovie.PosterPath,
		}
		movies = append(movies, movie)
	}

	return c.JSON(movies)
}
