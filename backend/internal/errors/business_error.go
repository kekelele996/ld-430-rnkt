package errors

import "fmt"

// BusinessError is the unified business error type carrying an error code and message.
type BusinessError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface.
func (e *BusinessError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("business error code=%d message=%s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("business error code=%d message=%s", e.Code, e.Message)
}

// Unwrap returns the wrapped cause for errors.Is / errors.As chains.
func (e *BusinessError) Unwrap() error { return e.Err }

// NewBusinessError creates a BusinessError with the given code and message.
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}

// WrapBusinessError creates a BusinessError wrapping an underlying error.
func WrapBusinessError(code int, message string, err error) *BusinessError {
	return &BusinessError{Code: code, Message: message, Err: err}
}
