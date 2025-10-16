package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/database"
	"gorm.io/gorm"
)

func BodyParser[T any](c *fiber.Ctx, dto *T) error {
	newReader := bytes.NewReader(c.Body())
	decoder := json.NewDecoder(newReader)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dto); err != nil {
		return errors.New("invalid or extra fields in request body")
	}
	return nil
}

func FindModelById[U any, T any](id U, model *T, modelTag string) error {
	if err := database.Database.Db.First(model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s does not exist", modelTag)
		}
		return errors.New("database error")
	}
	return nil
}
