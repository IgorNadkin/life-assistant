package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect postgres:", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db
}
