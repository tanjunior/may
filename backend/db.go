package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func db() *gorm.DB {

	// Load from environment variables
	dsn := "postgresql://neondb_owner:npg_8Bev2qcJCoRU@ep-mute-cake-a1wpo6dc-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to Railway PostgreSQL database!")

	return db

	// ctx := context.Background()

	// // Migrate the schema
	// db.AutoMigrate(&Product{})
}
