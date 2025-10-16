package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/services"
)

func SetupUserRoutes(api fiber.Router) {
	// Users endpoint
	user := api.Group("/users")
	user.Post("/", services.CreateUser)
	user.Get("/", services.GetUsers)
	user.Get("/:id", services.GetUser)
	user.Patch("/:id", services.UpdateUser)
	user.Delete("/:id", services.DeleteUser)
}
