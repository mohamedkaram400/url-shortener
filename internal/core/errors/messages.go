package domainerrors

import (
	"errors"
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

