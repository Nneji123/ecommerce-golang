package db

import (
	"github.com/nneji123/ecommerce-golang/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB *gorm.DB
)

// Connect to database
func Connect() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}
	dsn := cfg.PostgresDSN

	var db *gorm.DB
	for attempt := 1; attempt <= 3; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Error connecting to database (attempt %d): %s", attempt, err)
			// Wait for a moment before retrying
			time.Sleep(5 * time.Second)
			continue
		}
		// Successfully connected
		break
	}
	if err != nil {
		return nil, err
	}

	// Perform auto migration for multiple tables
	err = db.AutoMigrate(&LeadList{}, &Lead{})
	if err != nil {
		log.Fatalf("Error performing auto migration: %s", err)
	}

	DB = db

	return db, nil
}

// Close database connection
func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
