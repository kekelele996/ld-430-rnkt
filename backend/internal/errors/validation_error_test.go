package errors

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type sampleDTO struct {
	Email string `validate:"required,email"`
	Age   int    `validate:"required,min=18"`
}

func TestNewValidationError(t *testing.T) {
	t.Parallel()
	validate := validator.New()
	err := NewValidationError(validate.Struct(sampleDTO{}))
	if err == nil {
		t.Fatal("NewValidationError() returned nil")
	}
	if len(err.FieldErrors) == 0 {
		t.Fatalf("FieldErrors empty, got %#v", err.FieldErrors)
	}
	if err.Error() == "" {
		t.Fatal("Error() returned empty string")
	}
}

func TestNewValidationErrorNonValidatorError(t *testing.T) {
	t.Parallel()
	err := NewValidationError(&BusinessError{Code: 40000, Message: "bad"})
	if len(err.FieldErrors) != 1 || err.FieldErrors[0].Field != "body" {
		t.Fatalf("unexpected FieldErrors: %#v", err.FieldErrors)
	}
}

func TestBusinessErrorUnwrap(t *testing.T) {
	t.Parallel()
	cause := &BusinessError{Code: 40400, Message: "not found"}
	wrapped := WrapBusinessError(50000, "wrapper", cause)
	if wrapped.Unwrap() != cause {
		t.Fatalf("Unwrap() = %#v, want cause", wrapped.Unwrap())
	}
}
