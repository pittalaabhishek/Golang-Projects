package storage

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func (t Task) String() string {
	status := "incomplete"
	if t.Completed {
		status = "completed"
	}
	return fmt.Sprintf("Task{ID: %d, Description: %s, Status: %s}", t.ID, t.Description, status)
}

func (t Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t Task) IsCompletedToday() bool {
	if !t.Completed || t.CompletedAt == nil {
		return false
	}

	now := time.Now()
	completedDate := t.CompletedAt.Truncate(24 * time.Hour)
	todayDate := now.Truncate(24 * time.Hour)

	return completedDate.Equal(todayDate)
}