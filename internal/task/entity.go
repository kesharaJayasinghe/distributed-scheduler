package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID			uuid.UUID `json:"id"`
	Status		string    `json:"status"`
	Payload		[]byte    `json:"payload"`
	DueAt		time.Time `json:"due_at"`
	CreatedAt 	time.Time `json:"created_at"`
}