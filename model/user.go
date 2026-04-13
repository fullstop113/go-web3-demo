package model

import "gorm.io/gorm"

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    int64          `json:"-"`
	UpdatedAt    int64          `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"unique;not null" json:"username"`
	Email        string         `gorm:"unique;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
}