package db

import (
	"github.com/nneji123/ecommerce-golang/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB *gorm.DB
)

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
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}

	DB = db

	return db, nil
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
