package postgres

import (
	"backend/internal/domain/user"

	"github.com/jmoiron/sqlx"
)

type UserActionRepo struct {
	db *sqlx.DB
}

func NewUserActionRepo(db *sqlx.DB) *UserActionRepo {
	return &UserActionRepo{db: db}
}

func (r *UserActionRepo) Create(action *user.UserAction) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO user_actions (user_state_id, node_id, action, organization, deadline, completed) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		action.UserStateID, action.NodeID, action.Action, action.Organization, action.Deadline, action.Completed,
	).Scan(&id)

	return id, err
}

func (r *UserActionRepo) GetByUserState(userStateID int64) ([]user.UserAction, error) {
	var actions []user.UserAction
	err := r.db.Select(&actions,
		`SELECT id, user_state_id, node_id, action, organization, deadline, completed, completed_at, created_at
		 FROM user_actions WHERE user_state_id=$1 ORDER BY created_at DESC`,
		userStateID,
	)
	return actions, err
}

func (r *UserActionRepo) Update(action *user.UserAction) error {
	_, err := r.db.Exec(
		`UPDATE user_actions SET action=$1, organization=$2, deadline=$3, completed=$4, completed_at=$5 WHERE id=$6`,
		action.Action, action.Organization, action.Deadline, action.Completed, action.CompletedAt, action.ID,
	)
	return err
}
