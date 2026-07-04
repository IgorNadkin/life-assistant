package node

type Type string

const (
	Event    Type = "EVENT"
	Question Type = "QUESTION"
	Action   Type = "ACTION"
	Rule     Type = "RULE"
)

type Node struct {
	ID         int64
	ScenarioID int64
	Type       Type
	Title      string
	Content    string
	Meta       map[string]any
}
