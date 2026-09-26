package middleware

import (
	"github.com/assethub/assethub/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateJSON binds and validates a JSON body into dst.
func ValidateJSON(c *gin.Context, dst interface{}) error {
	if err := c.ShouldBindJSON(dst); err != nil {
		return errors.NewValidationError(err)
	}
	if err := validate.Struct(dst); err != nil {
		return errors.NewValidationError(err)
	}
	return nil
}

// ValidateForm binds and validates a multipart/form body into dst.
func ValidateForm(c *gin.Context, dst interface{}) error {
	if err := c.ShouldBind(dst); err != nil {
		return errors.NewValidationError(err)
	}
	if err := validate.Struct(dst); err != nil {
		return errors.NewValidationError(err)
	}
	return nil
}
