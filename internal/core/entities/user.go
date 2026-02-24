package entities

import "time"


type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName  string    `json:"username" gorm:"column:username;type:varchar(50);not null"`
	FirstName string    `json:"first_name" gorm:"column:first_name;type:varchar(50);not null"`
	LastName  string    `json:"last_name" gorm:"column:last_name;type:varchar(50);not null;uniqueIndex"`
	Email     string    `json:"email" gorm:"column:email;type:varchar(100);not null;uniqueIndex"`
	Password     string `json:"password" gorm:"column:password;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}