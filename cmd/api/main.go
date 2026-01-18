package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	log.Printf("Starting application in %s mode on port %s", cfg.Environment, cfg.Port)

	// Initialize database connection
	dbPool, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Initialize repository
	taskRepo := task.NewRepository(dbPool)

	// Define handler
	http.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		type RequestPayload struct {
			DueAt	time.Time       `json:"due_at"`
			Payload json.RawMessage `json:"payload"`
		}
		var req RequestPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Create task entity
		newTask := &task.Task{
			ID:			uuid.New(),
			Status:		"PENDING",
			Payload:	req.Payload,
			DueAt:		req.DueAt,
			CreatedAt:	time.Now(),
		}

		// Save to DB
		if err := taskRepo.Create(r.Context(), newTask); err != nil {
			log.Printf("Error creating task: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Respond
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"id": newTask.ID.String(),
			"status": "created",
		})
		
	})

	// Start server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
