package postgres

import (
	"backend/internal/domain/user"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(u *user.User) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (email, phone) VALUES ($1, $2) RETURNING id`,
		u.Email, u.Phone,
	).Scan(&id)

	return id, err
}

func (r *UserRepo) GetByID(id int64) (*user.User, error) {
	var u user.User
	err := r.db.Get(&u,
		`SELECT id, email, phone, created_at FROM users WHERE id=$1`,
		id,
	)
	return &u, err
}

func (r *UserRepo) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Get(&u,
		`SELECT id, email, phone, created_at FROM users WHERE email=$1`,
		email,
	)
	return &u, err
}
