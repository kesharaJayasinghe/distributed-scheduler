package task

import (
	"context"
	"fmt"

	// "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handles database operations for tasks
type Repository struct {
	db *pgxpool.Pool
}

// Create new instance
func NewRepository (db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Insert a new task into database
func (r *Repository) Create(ctx context.Context, t *Task) error {
	query := `
		INSERT INTO tasks (id, status, payload, due_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	// Execute the query using the connection pool
	_, err := r.db.Exec(ctx, query, t.ID, t.Status, t.Payload, t.DueAt, t.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	return nil
}