package main

import (
	"log"

	"github.com/kesharaJayasinghe/distributed-scheduler/internal/config"
	"github.com/kesharaJayasinghe/distributed-scheduler/internal/db"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting application in %s mode on port %s", cfg.Environment, cfg.Port)

	// Initialize database connection
	pool, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Successfully connected to database")
	log.Println("Application started successfully")
}
