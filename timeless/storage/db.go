package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./timeless.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("connected to db")

	DB = db
	return db, nil
}
