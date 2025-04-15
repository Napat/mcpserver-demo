package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator is a structure for using the validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate checks the validity of data
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// RegisterValidator registers the validator with Echo
func RegisterValidator(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
}
