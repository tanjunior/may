package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func db() *gorm.DB {

	// Load DSN from environment variables to avoid committing secrets.
	// Supports `DATABASE_URL` (common) or `POSTGRES_DSN` (as used in README).
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_DSN")
	}
	if dsn == "" {
		log.Fatal("DATABASE_URL or POSTGRES_DSN is not set")
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database!")

	return db

	// ctx := context.Background()

	// // Migrate the schema
	// db.AutoMigrate(&Product{})
}
