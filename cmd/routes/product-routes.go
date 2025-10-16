package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/services"
)

func SetupProductRoutes(api fiber.Router) {
	// Products endpoint
	product := api.Group("/products")
	product.Post("/", services.CreateProduct)
	product.Get("/", services.GetProducts)
	product.Get("/:id", services.GetProduct)
	product.Patch("/:id", services.UpdateProduct)
	product.Delete("/:id", services.DeleteProduct)
}
