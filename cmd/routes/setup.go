package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/services"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	// Welcome
	api.Get("/", services.WelcomeRoute)

	// User routes
	SetupUserRoutes(api)

	// Product routes
	SetupProductRoutes(api)

	// Order routes
	SetupOrderRoutes(api)
}
