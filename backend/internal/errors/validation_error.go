package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError aggregates field-level validation failures.
type ValidationError struct {
	FieldErrors []FieldError
}

// FieldError describes a single invalid field.
type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	parts := make([]string, 0, len(e.FieldErrors))
	for _, fe := range e.FieldErrors {
		parts = append(parts, fmt.Sprintf("%s: %s", fe.Field, fe.Tag))
	}
	return "validation failed: " + strings.Join(parts, ", ")
}

// NewValidationError converts validator errors into a ValidationError.
func NewValidationError(err error) *ValidationError {
	ve := &ValidationError{}
	var validationErrors validator.ValidationErrors
	if ok := asValidatorErrors(err, &validationErrors); ok {
		for _, fe := range validationErrors {
			ve.FieldErrors = append(ve.FieldErrors, FieldError{
				Field: fe.Field(),
				Tag:   fe.Tag(),
				Value: fe.Param(),
			})
		}
	} else {
		ve.FieldErrors = append(ve.FieldErrors, FieldError{Field: "body", Tag: "invalid", Value: err.Error()})
	}
	return ve
}

func asValidatorErrors(err error, target *validator.ValidationErrors) bool {
	if v, ok := err.(validator.ValidationErrors); ok {
		*target = v
		return true
	}
	return false
}
