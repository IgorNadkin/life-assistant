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
		`INSERT INTO scenario (title, description, category, start_node_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		s.Title, s.Description, s.Category, s.StartNodeID,
	).Scan(&id)

	return id, err
}

func (r *ScenarioRepo) Get(id int64) (*scenario.Scenario, error) {
	var s scenario.Scenario

	err := r.db.Get(&s,
		`SELECT id, title, description, category, start_node_id, created_at
		 FROM scenario WHERE id=$1`,
		id,
	)

	return &s, err
}

func (r *ScenarioRepo) Update(s *scenario.Scenario) error {
	_, err := r.db.Exec(
		`UPDATE scenario SET title=$1, description=$2, category=$3, start_node_id=$4 WHERE id=$5`,
		s.Title, s.Description, s.Category, s.StartNodeID, s.ID,
	)
	return err
}

func (r *ScenarioRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM scenario WHERE id=$1`, id)
	return err
}

func (r *ScenarioRepo) List() ([]scenario.Scenario, error) {
	var scenarios []scenario.Scenario
	err := r.db.Select(&scenarios,
		`SELECT id, title, description, category, start_node_id, created_at
		 FROM scenario ORDER BY created_at DESC`,
	)
	return scenarios, err
}
