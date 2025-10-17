package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/database"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/models"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/utils"
	"gorm.io/gorm"
)

func CreateProduct(c *fiber.Ctx) error {
	db := database.Database.Db
	var dto models.CreateProductDTO

	if err := utils.BodyParseAndValidate(c, &dto); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusBadRequest)
	}

	// Check if product already exists
	var existingProduct models.Product
	if err := db.Where("serial_number = ?", dto.SerialNumber).First(&existingProduct).Error; err == nil {
		// Found existing product with serial number
		return utils.ErrorResponse(c, "Product already exists", fiber.StatusConflict)
	} else if err != gorm.ErrRecordNotFound {
		// Some other DB error
		return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
	}

	product := models.Product{
		Name:         dto.Name,
		SerialNumber: dto.SerialNumber,
		Price:        dto.Price,
	}

	if err := db.Create(&product).Error; err != nil {
		return utils.ErrorResponse(c, "Could not create product", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "Product created successfully", models.CreateResponseProduct(product), fiber.StatusCreated)
}

func GetProducts(c *fiber.Ctx) error {
	db := database.Database.Db
	products := []models.Product{}

	if err := db.Preload("Orders.User").Find(&products).Error; err != nil {
		return utils.ErrorResponse(c, "Could not find products", fiber.StatusInternalServerError)
	}

	responseProducts := make([]models.ProductSerializer, len(products))

	for i, product := range products {
		responseProducts[i] = models.CreateResponseProduct(product)
	}

	return utils.SuccessResponse(c, "Products retrieved successfully", responseProducts)
}

func GetProduct(c *fiber.Ctx) error {
	db := database.Database.Db
	productId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid product ID", fiber.StatusBadRequest)
	}

	var product models.Product
	if err := utils.FindModelById(productId, &product, "product"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if err := db.Preload("Orders.User").First(&product).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to load products", fiber.StatusNotFound)
	}

	return utils.SuccessResponse(c, "Product retrieved successfully", models.CreateResponseProduct(product))
}

func UpdateProduct(c *fiber.Ctx) error {
	db := database.Database.Db
	productId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid product ID", fiber.StatusBadRequest)
	}

	var dto models.UpdateProductDTO
	if err := utils.BodyParseAndValidate(c, &dto); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusBadRequest)
	}

	var product models.Product
	if err := utils.FindModelById(productId, &product, "product"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if dto.Name != "" {
		product.Name = dto.Name
	}

	if dto.Price != 0 {
		product.Price = dto.Price
	}

	if dto.SerialNumber != "" {
		// Check if another product with the same serial number exists
		var existingProduct models.Product
		if err := db.Where("serial_number = ? AND id != ?", dto.SerialNumber, product.ID).First(&existingProduct).Error; err == nil {
			return utils.ErrorResponse(c, "Serial number is already taken by another product", fiber.StatusConflict)
		} else if err != gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
		}
		product.SerialNumber = dto.SerialNumber
	}

	if err := db.Save(&product).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to save products", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "Product updated successfully", models.CreateResponseProduct(product))
}

func DeleteProduct(c *fiber.Ctx) error {
	db := database.Database.Db
	productId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid product ID", fiber.StatusBadRequest)
	}

	var product models.Product
	if err := utils.FindModelById(productId, &product, "product"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if err := db.Delete(&product).Error; err != nil {
		return utils.ErrorResponse(c, "Failed to delete product", fiber.StatusBadRequest)
	}

	return utils.SuccessResponse(c, "Product deleted successfully", nil, fiber.StatusNoContent)
}
