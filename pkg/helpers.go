package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"math/big"
	"strings"
)

// Auth

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


// Url

const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateShortCode(shortCodeLenght int) (string, error) {
	shortCode := make([]byte, shortCodeLenght)	

	charsetLen := big.NewInt(int64(len(characters)))

	for i := 0; i < shortCodeLenght; i ++ {
		// Generate a random index within the charset length using crypto/rand
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		// Append the character at the random index to the result
		shortCode[i] = characters[int(randomIndex.Uint64())]
	}

	return string(shortCode), nil
}




// DB

func RunMigrations(dbURL string) {
	m, err := migrate.New(
		"file://./db/migrations", 
		dbURL,
	)
	if err != nil {
		log.Fatal("❌ Failed to create migrate instance:", err)
	}

	err = m.Up() // runs all pending migrations
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	log.Println("✅ Migrations ran successfully!")
}