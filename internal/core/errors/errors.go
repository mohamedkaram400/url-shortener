package domainerrors

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)



var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound        = errors.New("resource not found")
	ErrUserNotFound        = errors.New("user not found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")

	ErrShortCodeExists   = errors.New("custom alias already taken")
	ErrLinkExpired     = errors.New("link expired")
	ErrURLInactive       = errors.New("URL not active")
	ErrRefreshTokenMissing 		 = errors.New("refresh token missing")
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

	// 1️⃣ DTO validation errors (binding + validator)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {

			// Get JSON field name instead of struct field name
			field := fe.Field()
			if fe.StructField() != "" {
				field = getJSONFieldName(fe)
			}

			field = strings.ToLower(field)

			switch fe.Tag() {
			case "required":
				errorsMap[field] = "this field is required"
			case "email":
				errorsMap[field] = "invalid email format"
			case "min":
				errorsMap[field] = "value is too short"
			case "max":
				errorsMap[field] = "value is too long"
			default:
				errorsMap[field] = "invalid value"
			}
		}
		return errorsMap
	}

	// 2️⃣ Custom FieldError
	if fe, ok := err.(*FieldError); ok {
		errorsMap[strings.ToLower(fe.Field)] = fe.Message
		return errorsMap
	}

	// 3️⃣ Multi-field errors
	if mfe, ok := err.(*MultiFieldError); ok {
		for _, fe := range mfe.Errors {
			errorsMap[strings.ToLower(fe.Field)] = fe.Message
		}
		return errorsMap
	}

	return nil
}


func getJSONFieldName(fe validator.FieldError) string {
	field, _ := reflect.TypeOf(fe.StructNamespace()).FieldByName(fe.StructField())
	if tag, ok := field.Tag.Lookup("json"); ok {
		name := strings.Split(tag, ",")[0]
		return name
	}
	return fe.Field()
}