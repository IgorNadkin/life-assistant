package postgres

import (
	"backend/internal/domain/user"

	"github.com/jmoiron/sqlx"
)

type UserStateRepo struct {
	db *sqlx.DB
}

func NewUserStateRepo(db *sqlx.DB) *UserStateRepo {
	return &UserStateRepo{db: db}
}

func (r *UserStateRepo) Create(state *user.UserState) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO user_states (user_id, scenario_id, current_node_id, status, completed_steps) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		state.UserID, state.ScenarioID, state.CurrentNodeID, state.Status, state.CompletedSteps,
	).Scan(&id)

	return id, err
}

func (r *UserStateRepo) GetByID(id int64) (*user.UserState, error) {
	var state user.UserState
	err := r.db.Get(&state,
		`SELECT id, user_id, scenario_id, current_node_id, status, completed_steps, created_at, updated_at
		 FROM user_states WHERE id=$1`,
		id,
	)
	return &state, err
}

func (r *UserStateRepo) GetByUser(userID int64) (*user.UserState, error) {
	var state user.UserState
	err := r.db.Get(&state,
		`SELECT id, user_id, scenario_id, current_node_id, status, completed_steps, created_at, updated_at
		 FROM user_states WHERE user_id=$1 ORDER BY updated_at DESC LIMIT 1`,
		userID,
	)
	return &state, err
}

func (r *UserStateRepo) GetByUserAndScenario(userID, scenarioID int64) (*user.UserState, error) {
	var state user.UserState
	err := r.db.Get(&state,
		`SELECT id, user_id, scenario_id, current_node_id, status, completed_steps, created_at, updated_at
		 FROM user_states WHERE user_id=$1 AND scenario_id=$2 ORDER BY updated_at DESC LIMIT 1`,
		userID, scenarioID,
	)
	return &state, err
}

func (r *UserStateRepo) Update(state *user.UserState) error {
	_, err := r.db.Exec(
		`UPDATE user_states SET current_node_id=$1, status=$2, completed_steps=$3, updated_at=NOW() WHERE id=$4`,
		state.CurrentNodeID, state.Status, state.CompletedSteps, state.ID,
	)
	return err
}
