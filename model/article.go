package model

import "gorm.io/gorm"

type Article struct {
	ID uint `gorm:"primaryKey" json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeleteAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
	Title string `gorm:"not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	Status string `gorm:"not null;default:draft" json:"status"`
	AuthorID uint `gorm:"index;not null" json:"author_id"`
}
