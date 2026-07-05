package user

import "time"

type User struct {
	ID        int64     `db:"id"`
	Email     *string   `db:"email"`
	Phone     *string   `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}
