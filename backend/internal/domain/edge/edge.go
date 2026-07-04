package edge

type Edge struct {
	ID         int64
	FromNodeID int64
	ToNodeID   int64
	Condition  map[string]any
	Priority   int
}
