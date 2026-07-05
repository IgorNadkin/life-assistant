package graph

type Edge struct {
	ID        int64   `db:"id"`
	FromNode  int64   `db:"from_node"`
	ToNode    int64   `db:"to_node"`
	Condition string  `db:"condition"`
	Logic     *string `db:"logic"`
}
