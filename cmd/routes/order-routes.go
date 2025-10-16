package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/services"
)

func SetupOrderRoutes(api fiber.Router) {
	// Order endpoint
	order := api.Group("/orders")
	order.Post("/", services.CreateOrder)
	order.Get("/", services.GetOrders)
	order.Get("/:id", services.GetOrder)
	order.Patch("/:id", services.UpdateOrder)
	order.Delete("/:id", services.DeleteOrder)
}
