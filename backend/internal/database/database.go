package database

import (
	"log"

	"github.com/johnson7543/chatsheet-assessment/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(dbPath string) error {
	var err error

	// Open database connection
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.Println("Database connection established")

	// Run migrations
	if err := RunMigrations(); err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

// RunMigrations runs all database migrations
func RunMigrations() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.LinkedAccount{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
