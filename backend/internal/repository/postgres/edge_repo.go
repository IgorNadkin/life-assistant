package postgres

import (
	"backend/internal/domain/graph"

	"github.com/jmoiron/sqlx"
)

type EdgeRepo struct {
	db *sqlx.DB
}

func NewEdgeRepo(db *sqlx.DB) *EdgeRepo {
	return &EdgeRepo{db: db}
}

func (r *EdgeRepo) Create(edge *graph.Edge) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO graph_edge (from_node, to_node, condition, logic) VALUES ($1, $2, $3, $4) RETURNING id`,
		edge.FromNode, edge.ToNode, edge.Condition, edge.Logic,
	).Scan(&id)

	return id, err
}

func (r *EdgeRepo) GetByID(id int64) (*graph.Edge, error) {
	var edge graph.Edge
	err := r.db.Get(&edge,
		`SELECT id, from_node, to_node, condition, logic FROM graph_edge WHERE id=$1`,
		id,
	)
	return &edge, err
}

func (r *EdgeRepo) GetByScenario(scenarioID int64) ([]graph.Edge, error) {
	var edges []graph.Edge
	err := r.db.Select(&edges,
		`SELECT e.id, e.from_node, e.to_node, e.condition, e.logic FROM graph_edge e
		 INNER JOIN graph_node n ON e.from_node = n.id
		 WHERE n.scenario_id=$1`,
		scenarioID,
	)
	return edges, err
}

func (r *EdgeRepo) Update(edge *graph.Edge) error {
	_, err := r.db.Exec(
		`UPDATE graph_edge SET from_node=$1, to_node=$2, condition=$3, logic=$4 WHERE id=$5`,
		edge.FromNode, edge.ToNode, edge.Condition, edge.Logic, edge.ID,
	)
	return err
}

func (r *EdgeRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM graph_edge WHERE id=$1`, id)
	return err
}
