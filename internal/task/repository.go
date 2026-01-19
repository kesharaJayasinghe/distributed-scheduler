package task

import (
	"context"
	"fmt"
	"time"

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

// Find tasks that are due and mark as 'RUNNING'
func (r *Repository) ClaimDueTasks(ctx context.Context) ([]Task, error) {

	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)	// Defer rollback in case of error

	query := `
		UPDATE tasks
		SET status = 'RUNNING', picked_at = NOW()
		WHERE id IN (
			SELECT id
			FROM tasks
			WHERE status = 'PENDING' AND due_at <= NOW()
			ORDER BY due_at ASC
			LIMIT 10
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, status, payload, due_at, created_at
	`

	rows, err := tx.Query(ctx, query)
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

	// Commit the transaction to finalize the 'RUNNING' state
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update task status
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, "UPDATE tasks SET status = $1 WHERE id = $2", status, id)
	return err
}

// Find task stuck in RUNNING state for too long and reset to PENDING
func (r *Repository) ResetZombieTasks(ctx context.Context, maxDuration time.Duration) (int64, error) {
	// Calculate the cutoff time. Any task started before this time is considered dead.
	cutoff := time.Now().Add(-maxDuration)

	query := `
		UPDATE tasks
		SET status = 'PENDING', picked_at = NULL
		WHERE status = 'RUNNING'
		AND picked_at < $1
	`

	tag, err := r.db.Exec(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}