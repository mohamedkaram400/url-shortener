package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type CustomClaims struct {
	UserID    int64      `json:"sub"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

type Session struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Device       string     `json:"device" gorm:"column:device;type:varchar(255);not null"`
	RefreshToken string     `json:"refresh_token" gorm:"column:refresh_token;type:text;not null"`
	UserID       uint       `json:"user_id" gorm:"column:user_id;not null;index"`
	IpAddress    string     `json:"ip_address" gorm:"column:ip_address;type:varchar(50)"`
	ExpiresAt    time.Time  `json:"expires_at" gorm:"column:expires_at;not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relation
	User         User       `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
