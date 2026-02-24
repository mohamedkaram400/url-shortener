package pkg

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string) {
	// The cost parameter controls how slow the hash is.
	// bcrypt automatically generates a random salt.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
        log.Fatal("❌ Error generating hashed password:", err)
	}
	return string(bytes)
}


func CheckPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}