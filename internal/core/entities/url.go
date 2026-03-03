package entities

import "time"


type Url struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	OriginalURL string     `json:"original_url" gorm:"column:original_url;type:varchar(255);not null"`
	ShortCode   string     `json:"short_code" gorm:"column:short_code;type:varchar(50);not null;uniqueIndex"`
	Status      string     `json:"status" gorm:"column:status;type:varchar(20);not null;default:'Active'"`
	UserID      uint       `json:"user_id" gorm:"column:user_id;not null;index"`
	ClickCount  int        `json:"click_count" gorm:"column:click_count;default:0"`
	CustomAlias *string    `json:"custom_alias" gorm:"column:custom_alias;varchar(20)"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`

	User         User       `json:"user" gorm:"foreignKey:UserID;references:ID"`
}