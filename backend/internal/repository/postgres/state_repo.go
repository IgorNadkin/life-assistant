package postgres

import "github.com/jmoiron/sqlx"

type StateRepo struct {
	db *sqlx.DB
}

func NewStateRepo(db *sqlx.DB) *StateRepo {
	return &StateRepo{db: db}
}

func (r *StateRepo) Get(userID string) (int64, error) {
	var nodeID int64

	err := r.db.Get(&nodeID,
		`SELECT current_node_id FROM user_state WHERE user_id=$1`,
		userID,
	)

	return nodeID, err
}

func (r *StateRepo) Save(userID string, nodeID int64) error {
	_, err := r.db.Exec(`
		INSERT INTO user_state (user_id, current_node_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET current_node_id=$2, updated_at=NOW()
	`, userID, nodeID)

	return err
}
