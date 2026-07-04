package postgres

import (
	"backend/internal/domain/scenario"

	"github.com/jmoiron/sqlx"
)

type ScenarioRepo struct {
	db *sqlx.DB
}

func NewScenarioRepo(db *sqlx.DB) *ScenarioRepo {
	return &ScenarioRepo{db: db}
}

func (r *ScenarioRepo) Create(s *scenario.Scenario) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO scenario (title, description, start_node_id) VALUES ($1, $2, $3) RETURNING id`,
		s.Title, s.Description, s.StartNodeID,
	).Scan(&id)

	return id, err
}

func (r *ScenarioRepo) Get(id int64) (*scenario.Scenario, error) {
	var s scenario.Scenario

	err := r.db.Get(&s,
		`SELECT id, title, description, start_node_id
		 FROM scenario WHERE id=$1`,
		id,
	)

	return &s, err
}
