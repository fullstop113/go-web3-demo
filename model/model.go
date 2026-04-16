package model

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	// 初始化数据库
	if err := DB.AutoMigrate(&User{}, &Article{}); err != nil {
		log.Fatal(err)
	}
}

const Name = "model"
