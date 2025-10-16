package utils

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data any, status ...int) error {
	code := fiber.StatusOK
	if len(status) > 0 {
		code = status[0]
	}

	return c.Status(code).JSON(ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, message string, status int, data ...any) error {
	var responseData any = nil

	if len(data) > 0 {
		responseData = data[0]
	}

	return c.Status(status).JSON(ApiResponse{
		Success: false,
		Message: message,
		Data:    responseData,
	})
}

func SuccessResponseWithStatus(c *fiber.Ctx, message string, data any, status int) error {
	return c.Status(status).JSON(ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
