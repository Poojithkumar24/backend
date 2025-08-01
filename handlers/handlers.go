package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	db "github.com/Poojithkumar24/backend/db"
	"github.com/Poojithkumar24/backend/types"
	"github.com/t3rm1n4l/go-mega"

	"github.com/gofiber/fiber/v2"
)

func AddMovie(c *fiber.Ctx) error {

	log.Println("add mocie called ")
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

func UploadHandler(c *fiber.Ctx) error {
	// Get file from form-data (field name should be "file")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to get file: "+err.Error())
	}

	// Open file stream
	src, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open uploaded file: "+err.Error())
	}
	defer src.Close()

	// Save uploaded file to a temporary path
	tempDir := os.TempDir()
	localPath := filepath.Join(tempDir, fileHeader.Filename)

	dst, err := os.Create(localPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create temp file: "+err.Error())
	}
	defer dst.Close()
	io.Copy(dst, src)

	// Upload to MEGA
	link, key, fileName, err := UploadToMega(localPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload to MEGA: "+err.Error())
	}

	// Optional: Clean up the temp file
	_ = os.Remove(localPath)

	// Return JSON response
	return c.JSON(fiber.Map{
		"fileName":      fileName,
		"megaLink":      link,
		"decryptionKey": key,
	})
}

func UploadToMega(filePath string) (shareLink, decryptionKey, fileName string, err error) {
	filePath = strings.TrimSpace(filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		return "", "", "", fmt.Errorf("invalid file path: %v", err)
	}

	client := mega.New()
	if err := client.Login(email, password); err != nil {
		return "", "", "", fmt.Errorf("failed to login to MEGA: %v", err)
	}

	root := client.FS.GetRoot()
	if root == nil {
		return "", "", "", fmt.Errorf("failed to get root folder")
	}

	progressChan := make(chan int, 100)
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
			done <- true
		}()
		for range progressChan {
			// Progress is optional here
		}
	}()

	fileName = filepath.Base(filePath)
	node, err := client.UploadFile(filePath, root, fileName, &progressChan)
	<-done

	if err != nil {
		return "", "", "", fmt.Errorf("upload failed: %v", err)
	}

	fullLink, err := client.Link(node, true) // include key
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate link: %v", err)
	}

	parts := strings.Split(fullLink, "#")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("failed to extract decryption key from link")
	}

	return fullLink, parts[1], fileName, nil
}

func MegaDownloadHandler(c *fiber.Ctx) error {
	var req types.DownloadRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if !strings.Contains(req.MegaURL, "#") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid MEGA URL")
	}

	// Execute megadl command to download the file
	cmd := exec.Command("megadl", req.MegaURL, "--path", "/tmp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("megadl error: %s", output)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to download from MEGA")
	}

	// Extract filename from output (e.g., "Downloaded my_file.pdf")
	var fileName string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Downloaded ") {
			fileName = strings.TrimPrefix(line, "Downloaded ")
			break
		}
	}

	if fileName == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse downloaded filename")
	}

	path := filepath.Join("/tmp", fileName)
	return c.Download(path)

}
