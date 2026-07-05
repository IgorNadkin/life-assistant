package handler

import (
	"encoding/json"
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
	if err := h.db.Ping(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
