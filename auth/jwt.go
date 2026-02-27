package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)


type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

func ValidateJWT(secretKey string, tokenString string) (*entities.CustomClaims, error) {
	claims := &entities.CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil || !token.Valid {
		fmt.Println("JWT parse error:", err)
		return nil, err
	}

	return claims, nil
}

func IsuueTokens(secretKey string, AccessTokenTime int, RefrashTokenTime int, userID uint) (string, string, error) {
	accessToken, err := GenerateAccessToken(secretKey, AccessTokenTime, userID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenerateRefrashToken(secretKey, userID, RefrashTokenTime)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateRefrashToken(secretKey string, userID uint, RefrashTokenTime int) (string, error) {
	return GenerateToken(secretKey, time.Duration(RefrashTokenTime) * 24 * time.Hour, userID, RefreshToken)
}

func GenerateAccessToken(secretKey string, AccessTokenTime int, userID uint) (string, error) {
	return GenerateToken(secretKey, time.Duration(AccessTokenTime) * time.Minute, userID, AccessToken)
}

func GenerateToken(secretKey string, duration time.Duration, userID uint, tokenType TokenType) (string, error) {
	jwtSecret := []byte(secretKey)
	
	claims := entities.CustomClaims{
		UserID:    int64(userID),
		TokenType: string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}