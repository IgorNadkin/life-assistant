package scenario

import "time"

type Scenario struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Category    *string   `db:"category"`
	StartNodeID int64     `db:"start_node_id"`
	CreatedAt   time.Time `db:"created_at"`
}
