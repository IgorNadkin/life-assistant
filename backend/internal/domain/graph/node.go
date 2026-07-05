package graph

import "time"

type NodeType string

const (
	NodeTypeQuestion NodeType = "question"
	NodeTypeInfo     NodeType = "info"
	NodeTypeAction   NodeType = "action"
	NodeTypeResult   NodeType = "result"
)

type Node struct {
	ID              int64     `db:"id"`
	ScenarioID      int64     `db:"scenario_id"`
	Type            NodeType  `db:"type"`
	Text            string    `db:"text"`
	ActionType      *string   `db:"action_type"`
	Organization    *string   `db:"organization"`
	Deadline        *string   `db:"deadline"`
	ReferenceLinks  *string   `db:"reference_links"` // JSON array stored as string
	Consequences    *string   `db:"consequences"`
	Order           int       `db:"order"`
	CreatedAt       time.Time `db:"created_at"`
}
