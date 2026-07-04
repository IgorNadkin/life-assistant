package postgres

import (
	"backend/internal/domain/graph"

	"github.com/jmoiron/sqlx"
)

type GraphRepo struct {
	db *sqlx.DB
}

func NewGraphRepo(db *sqlx.DB) *GraphRepo {
	return &GraphRepo{db: db}
}

func (r *GraphRepo) GetNodes() ([]graph.Node, error) {
	var nodes []graph.Node

	err := r.db.Select(&nodes,
		`SELECT id, type, text FROM graph_node`,
	)

	return nodes, err
}

func (r *GraphRepo) GetEdges() ([]graph.Edge, error) {
	var edges []graph.Edge

	err := r.db.Select(&edges,
		`SELECT from_node as "from", to_node as "to", condition FROM graph_edge`,
	)

	return edges, err
}
