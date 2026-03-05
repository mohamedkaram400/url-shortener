package domainerrors

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)



type FieldError struct {
	Field   string
	Message string
}

type MultiFieldError struct {
	Errors []*FieldError
}

func (e *FieldError) Error() string {
	return e.Message
}

func (e *MultiFieldError) Error() string {
	return "multiple validation errors"
}

func NewFieldError(field, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

func ValidationError(err error) map[string]string {

	errorsMap := make(map[string]string)

	// log.Print("errorsMap: ", errorsMap)

	// DTO validation errors
	if ve, ok := err.(validator.ValidationErrors); ok {

		for _, fe := range ve {

			field := toSnakeCase(fe.Field())

			errorsMap[field] = formatValidationError(fe)
		}
		log.Print("errorsMap1: ", errorsMap, "ValidationErrors: ", ve)

		return errorsMap
	}

	// Custom FieldError
	if fe, ok := err.(*FieldError); ok {
		errorsMap[fe.Field] = fe.Message

		log.Print("errorsMap1: ", errorsMap, "ValidationErrors: ", fe)
		
		return errorsMap
	}

	// Multi-field errors
	if mfe, ok := err.(*MultiFieldError); ok {
		for _, fe := range mfe.Errors {
			errorsMap[fe.Field] = fe.Message
		}
		log.Print("errorsMap1: ", errorsMap, "ValidationErrors: ", mfe)

		return errorsMap
	}
	
	errorsMap["error"] = err.Error()
	return errorsMap
}

func formatValidationError(fe validator.FieldError) string {

	switch fe.Tag() {

		case "required":
			return fe.Field() + " is required"

		case "email":
			return fe.Field() + " must be a valid email address"

		case "min":
			return fe.Field() + " must be at least " + fe.Param() + " characters"

		case "max":
			return fe.Field() + " must not exceed " + fe.Param() + " characters"

		case "len":
			return fe.Field() + " must be exactly " + fe.Param() + " characters"

		case "alphanum":
			return fe.Field() + " must contain only letters and numbers"

		default:
			return fe.Field() + " is invalid"
	}

}

func toSnakeCase(field string) string {
	return strings.ToLower(field)
}