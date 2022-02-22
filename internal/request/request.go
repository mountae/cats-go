// Package request qwe
package request

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/google/uuid"
)

// CustomValidator struct init
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate func provides validation
func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// CatID struct init
type CatID struct {
	ID   uuid.UUID `json:"id" bson:"id"`
	Name string    `json:"name" bson:"name"`
}
