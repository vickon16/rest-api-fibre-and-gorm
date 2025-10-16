package services

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/database"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/models"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/utils"

	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx) error {
	db := database.Database.Db
	var dto models.CreateUserDTO

	// Parse the dto
	if err := c.BodyParser(&dto); err != nil {
		return utils.ErrorResponse(c, "Invalid request body", fiber.StatusBadRequest)
	}

	// Validate the input
	errs, _ := utils.ValidateDto(dto)
	if errs != nil {
		return utils.ErrorResponse(c, "Validation errors occurred", fiber.StatusBadRequest, errs)
	}

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		// Found existing user
		return utils.ErrorResponse(c, "User already exists with this email", fiber.StatusConflict)
	} else if err != gorm.ErrRecordNotFound {
		// Some other DB error
		return utils.ErrorResponse(c, "Database error", fiber.StatusInternalServerError)
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return utils.ErrorResponse(c, "Failed to hash password", fiber.StatusInternalServerError)
	}

	user := models.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  hashedPassword,
	}

	if err := db.Create(&user).Error; err != nil {
		return utils.ErrorResponse(c, "Could not create user", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "User created successfully", models.CreateResponseUser(user), fiber.StatusCreated)
}

func GetUsers(c *fiber.Ctx) error {
	db := database.Database.Db
	users := []models.User{}

	db.Find(&users)
	responseUsers := []models.UserSerializer{}

	for _, user := range users {
		responseUser := models.CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", responseUsers)
}

func GetUser(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid user ID", fiber.StatusBadRequest)
	}

	var user models.User
	if err := utils.FindModelById(userId, &user, "user"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	return utils.SuccessResponse(c, "User retrieved successfully", models.CreateResponseUser(user))
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.Database.Db
	userId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid user ID", fiber.StatusBadRequest)
	}

	var dto models.UpdateUserDTO

	// Parse the dto
	if err := utils.BodyParser(c, &dto); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusBadRequest)
	}

	// Validate the input
	if errs, err := utils.ValidateDto(dto); err != nil {
		log.Println("Unexpected validation error:", err)
	} else if errs != nil {
		return utils.ErrorResponse(c, "Validation errors occurred", fiber.StatusBadRequest, errs)
	}

	var user models.User
	if err := utils.FindModelById(userId, &user, "user"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}

	if dto.LastName != "" {
		user.LastName = dto.LastName
	}

	if dto.Email != "" {
		// Ensure no other user has this email
		var existingUser models.User
		if err := db.Where("email = ? AND id != ?", dto.Email, user.ID).First(&existingUser).Error; err == nil {
			return utils.ErrorResponse(c, "Email is already taken by another user", fiber.StatusConflict)
		}

		user.Email = dto.Email
	}

	// if dto.Password != "" {
	// 	hashedPassword, err := utils.HashPassword(dto.Password)
	// 	if err != nil {
	// return utils.ErrorResponse(c, "Failed to hash password", fiber.StatusInternalServerError)
	// 	user.Password = hashedPassword
	// }

	if err := db.Save(&user).Error; err != nil {
		return utils.ErrorResponse(c, "Could not update user", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "User updated successfully", models.CreateResponseUser(user))
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.Database.Db
	userId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, "Invalid user ID", fiber.StatusBadRequest)
	}

	var user models.User
	if err := utils.FindModelById(userId, &user, "user"); err != nil {
		return utils.ErrorResponse(c, err.Error(), fiber.StatusNotFound)
	}

	if err := db.Delete(&user).Error; err != nil {
		return utils.ErrorResponse(c, "Could not delete user", fiber.StatusInternalServerError)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil, fiber.StatusNoContent)
}
