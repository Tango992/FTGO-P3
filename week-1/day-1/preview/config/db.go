package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	
	if err != nil {
		log.Fatal(err)
	}
	// AutoMigrate(db)
	return db
}

// func AutoMigrate(db *gorm.DB) {
// 	if err := db.AutoMigrate(&entity.User{}, &entity.Verification{}, &entity.Category{}, &entity.Brand{}, &entity.Catalog{}, &entity.Order{}, &entity.Payment{}); err != nil {
// 		log.Fatal(err)
// 	}
// }
