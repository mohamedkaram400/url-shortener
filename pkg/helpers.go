package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// The cost parameter controls how slow the hash is.
	// bcrypt automatically generates a random salt.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}


func HashToken(refreshToken string) (string) {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}

func CheckToken(hashed string, plainToken string) bool {
	hashedToken := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(hashedToken[:]) == hashed
}


func ExtractTokenFromHeader(authHeader string) (string, error) {

	if authHeader == "" {
		return "", errors.New("invalid authorization header")
	}

	// Split "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}
