package main

import (
	"context"
	"log"
	"time"

	"github.com/kesharaJayasinghe/distributed-scheduler/internal/config"
	"github.com/kesharaJayasinghe/distributed-scheduler/internal/db"
	"github.com/kesharaJayasinghe/distributed-scheduler/internal/task"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbPool, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	taskRepo := task.NewRepository(dbPool)

	log.Println("Scheduler started, polling for tasks...")

	// Polling loop
	for {
		// find due tasks
		ctx := context.Background()
		tasks, err := taskRepo.ListDueTasks(ctx)
		if err != nil {
			log.Printf("Error fetching tasks: %v", err)
			time.Sleep(5 * time.Second) // Backoff on error
			continue
		}

		if len(tasks) == 0 {
			time.Sleep(2 * time.Second) // No tasks, wait before polling again
			continue
		}

		// Process each task
		for _, t := range tasks {
			log.Printf("Found task %s. Executing...", t.ID)

			// Simulate task
			time.Sleep(500 * time.Millisecond)

			// Mark as done
			if err := taskRepo.UpdateStatus(ctx, t.ID, "COMPLETED"); err != nil {
				log.Printf("Failed to mark task %s as COMPLETED: %v", t.ID, err)
			} else {
				log.Printf("Task %s completed successfully.", t.ID)
			}
		}
	}
}