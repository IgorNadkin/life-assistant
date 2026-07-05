package postgres

import (
	"backend/internal/domain/graph"

	"github.com/jmoiron/sqlx"
)

type NodeRepo struct {
	db *sqlx.DB
}

func NewNodeRepo(db *sqlx.DB) *NodeRepo {
	return &NodeRepo{db: db}
}

func (r *NodeRepo) Create(node *graph.Node) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO graph_node (scenario_id, type, text, action_type, organization, deadline, reference_links, consequences, "order") 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		node.ScenarioID, node.Type, node.Text, node.ActionType, node.Organization, node.Deadline, node.ReferenceLinks, node.Consequences, node.Order,
	).Scan(&id)

	return id, err
}

func (r *NodeRepo) GetByID(id int64) (*graph.Node, error) {
	var node graph.Node
	err := r.db.Get(&node,
		`SELECT id, scenario_id, type, text, action_type, organization, deadline, reference_links, consequences, "order", created_at
		 FROM graph_node WHERE id=$1`,
		id,
	)
	return &node, err
}

func (r *NodeRepo) GetByScenario(scenarioID int64) ([]graph.Node, error) {
	var nodes []graph.Node
	err := r.db.Select(&nodes,
		`SELECT id, scenario_id, type, text, action_type, organization, deadline, reference_links, consequences, "order", created_at
		 FROM graph_node WHERE scenario_id=$1 ORDER BY "order"`,
		scenarioID,
	)
	return nodes, err
}

func (r *NodeRepo) Update(node *graph.Node) error {
	_, err := r.db.Exec(
		`UPDATE graph_node SET type=$1, text=$2, action_type=$3, organization=$4, deadline=$5, reference_links=$6, consequences=$7, "order"=$8 WHERE id=$9`,
		node.Type, node.Text, node.ActionType, node.Organization, node.Deadline, node.ReferenceLinks, node.Consequences, node.Order, node.ID,
	)
	return err
}

func (r *NodeRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM graph_node WHERE id=$1`, id)
	return err
}
