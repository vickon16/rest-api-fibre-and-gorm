package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/database"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/routes"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/utils"
)

func main() {
	_ = godotenv.Load()
	database.ConnectDb()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := err.Error()

			// Check if it's a Fiber error (like 404, 400, etc.)
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			// Log internal server errors for debugging
			if code == fiber.StatusInternalServerError {
				log.Printf("Internal Server Error: %v", err)
			}

			return utils.ErrorResponse(c, message, code)
		},
	})

	// Middleware to recover from panics, ensure the app never crashes
	app.Use(recover.New())

	// ✅ Enable CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://yourfrontenddomain.com", // specify allowed origins
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Setup all routes with /api prefix
	routes.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4000" // default fallback
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
