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

func ValidationError(err error) map[string]string {
    errors := make(map[string]string)

    if ve, ok := err.(validator.ValidationErrors); ok {
        for _, fe := range ve {
            errors[fe.Field()] = fe.Tag()
        }
    }

    return errors
}