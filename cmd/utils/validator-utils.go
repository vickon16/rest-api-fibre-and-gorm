package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ValidateDto(dto any) (map[string]string, error) {
	if err := Validate.Struct(dto); err != nil {

		// Check if the error is of type validator.ValidationErrors
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := make(map[string]string)

			for _, e := range errs {
				validationErrors[e.Field()] = fmt.Sprintf(
					"Field validation for '%s' failed on the '%s' rule", e.Field(), e.Tag(),
				)
			}
			return validationErrors, nil
		}

		// Return non-validation errors (e.g. bad struct)
		return nil, err
	}

	// No errors
	return nil, nil
}
