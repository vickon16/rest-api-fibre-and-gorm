package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/utils"
)

func WelcomeRoute(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, "Welcome to the Go Fiber and GORM REST API", nil)
}
