package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type DBHealthHandler struct {
	db *sqlx.DB
}

func NewDBHealthHandler(db *sqlx.DB) *DBHealthHandler {
	return &DBHealthHandler{db: db}
}

func (h *DBHealthHandler) Check(w http.ResponseWriter, r *http.Request) {

	err := h.db.Ping()
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"database":"disconnected"}`))
		return
	}

	w.Write([]byte(`{"database":"connected"}`))
}
