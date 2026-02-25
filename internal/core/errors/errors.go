package coreErrors

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type FieldError struct {
	Field   string
	Message string
}

func NewFieldError(field, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

func (e *FieldError) Error() string {
	return e.Message
}

type MultiFieldError struct {
	Errors []*FieldError
}

func (e *MultiFieldError) Error() string {
	return "multiple validation errors"
}

func ValidationError(err error) map[string]string {
	errorsMap := make(map[string]string)

	// DTO validation errors
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			errorsMap[fe.Field()] = fe.Tag()
		}
		return errorsMap
	}

	// Custom FieldError
	if fe, ok := err.(*FieldError); ok {
		errorsMap[fe.Field] = fe.Message
		return errorsMap
	}

	// Multi-field errors
	if mfe, ok := err.(*MultiFieldError); ok {
		for _, fe := range mfe.Errors {
			errorsMap[fe.Field] = fe.Message
		}
		return errorsMap
	}

	// fallback
	errorsMap["error"] = err.Error()
	return errorsMap
}