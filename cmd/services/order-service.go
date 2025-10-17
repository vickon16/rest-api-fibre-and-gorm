package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/database"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/models"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/utils"
	"gorm.io/gorm"
)

func CreateOrder(c *fiber.Ctx) error {
	db := database.Database.Db
	var dto models.CreateOrderDTO

	if err := utils.BodyParseAndValidate(c, &dto); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusBadRequest)
	}

	// Run both existence check in parallel
	channelNumbers := 2
	errChan := make(chan error, channelNumbers)
	var existingProduct models.Product
	var existingUser models.User

	go func() {
		errChan <- db.Select("id").First(&existingProduct, dto.ProductId).Error
	}()

	go func() {
		errChan <- db.Select("id").First(&existingUser, dto.UserId).Error
	}()

	// Wait for both go routines
	for range channelNumbers {
		if err := <-errChan; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrorResponse(c, "Product or User not found", fiber.StatusNotFound)
			}
			return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
		}
	}

	order := models.Order{
		ProductID: dto.ProductId,
		UserID:    dto.UserId,
	}

	// ✅ Use a transaction for atomicity
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Preload associated only after creation
		if err := tx.Preload("Product").Preload("User").First(&order, order.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return utils.ErrorResponse(c, "Could not create order", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "Order created successfully", models.CreateResponseOrder(order), fiber.StatusCreated)
}

func GetOrders(c *fiber.Ctx) error {
	db := database.Database.Db
	orders := []models.Order{}

	if err := db.Preload("Product").Preload("User").Find(&orders).Error; err != nil {
		return utils.ErrorResponse(c, "Could not find orders", fiber.StatusInternalServerError)
	}

	responseOrders := make([]models.OrderSerializer, len(orders))

	for i, order := range orders {
		responseOrders[i] = models.CreateResponseOrder(order)
	}

	return utils.SuccessResponse(c, "Orders retrieved successfully", responseOrders)
}

func GetOrder(c *fiber.Ctx) error {
	db := database.Database.Db
	orderId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid order ID", fiber.StatusBadRequest)
	}

	var order models.Order
	if err := utils.FindModelById(orderId, &order, "order"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if err := db.Preload("Product").Preload("User").First(&order).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to load order", fiber.StatusNotFound)
	}

	return utils.SuccessResponse(c, "Order retrieved successfully", models.CreateResponseOrder(order))
}

func UpdateOrder(c *fiber.Ctx) error {
	db := database.Database.Db

	orderId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid order ID", fiber.StatusBadRequest)
	}

	var dto models.UpdateOrderDTO
	if err := utils.BodyParseAndValidate(c, &dto); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusBadRequest)
	}

	// Fetch the order (only once)
	var order models.Order
	if err := db.First(&order, orderId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, "Order not found", fiber.StatusNotFound)
		}
		return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
	}

	// Optional: check if new ProductID exists
	if dto.ProductId != 0 && dto.ProductId != order.ProductID {
		var count int64
		if err := db.Model(&models.Product{}).Where("id = ?", dto.ProductId).Count(&count).Error; err != nil {
			return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
		}
		if count == 0 {
			return utils.ErrorResponse(c, "Product not found", fiber.StatusNotFound)
		}
		order.ProductID = dto.ProductId
	}

	// Optional: check if new UserID exists
	if dto.UserId != 0 && dto.UserId != order.UserID {
		var count int64
		if err := db.Model(&models.User{}).Where("id = ?", dto.UserId).Count(&count).Error; err != nil {
			return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
		}
		if count == 0 {
			return utils.ErrorResponse(c, "User not found", fiber.StatusNotFound)
		}
		order.UserID = dto.UserId
	}

	// Save changes
	if err := db.Save(&order).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to update order", fiber.StatusInternalServerError)
	}

	// Preload associations for response
	if err := db.Preload("Product").Preload("User").First(&order, order.ID).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to load updated order", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "Order updated successfully", models.CreateResponseOrder(order))
}

func DeleteOrder(c *fiber.Ctx) error {
	db := database.Database.Db
	orderId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid order ID", fiber.StatusBadRequest)
	}

	var order models.Order
	if err := utils.FindModelById(orderId, &order, "order"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if err := db.Delete(&order).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to delete order", fiber.StatusBadRequest)
	}

	return utils.SuccessResponse(c, "Order deleted successfully", nil, fiber.StatusNoContent)
}
