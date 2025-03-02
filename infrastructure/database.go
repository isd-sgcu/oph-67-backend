package infrastructure

import (
	"fmt"
	"github.com/isd-sgcu/oph-67-backend/config"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDatabase(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to the database")

	// Automatically migrate the schema, creating tables if they don't exist
	err = db.AutoMigrate(&domain.User{}, &domain.StudentTransaction{}) // Add your domain models here
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	return db
}
