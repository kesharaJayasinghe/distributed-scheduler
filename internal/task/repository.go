package task

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

// Find tasks that are ready to run
func (r *Repository) ListDueTasks(ctx context.Context) ([]Task, error) {
	// Limit to 10
	query := `
		SELECT id, status, payload, due_at, created_at
		FROM tasks
		WHERE status = 'PENDING' AND due_at <= NOW()
		ORDER BY due_at ASC
		LIMIT 10
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Status, &t.Payload, &t.DueAt, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Update task status
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, "UPDATE tasks SET status = $1 WHERE id = $2", status, id)
	return err
}