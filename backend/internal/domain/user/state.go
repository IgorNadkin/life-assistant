package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type StateStatus string

const (
	StatusInProgress StateStatus = "in_progress"
	StatusCompleted  StateStatus = "completed"
)

type UserState struct {
	ID              int64       `db:"id"`
	UserID          int64       `db:"user_id"`
	ScenarioID      int64       `db:"scenario_id"`
	CurrentNodeID   int64       `db:"current_node_id"`
	Status          StateStatus `db:"status"`
	CompletedSteps  Int64Array  `db:"completed_steps"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
}

// Int64Array for PostgreSQL array support
type Int64Array []int64

func (a Int64Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Int64Array) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion failed")
	}
	return json.Unmarshal(bytes, &a)
}

type UserAction struct {
	ID           int64     `db:"id"`
	UserStateID  int64     `db:"user_state_id"`
	NodeID       int64     `db:"node_id"`
	Action       string    `db:"action"`
	Organization string    `db:"organization"`
	Deadline     *time.Time `db:"deadline"`
	Completed    bool      `db:"completed"`
	CompletedAt  *time.Time `db:"completed_at"`
	CreatedAt    time.Time `db:"created_at"`
}
