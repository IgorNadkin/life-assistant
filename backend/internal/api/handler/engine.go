package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/repository/postgres"
	"backend/internal/service/engine"
)

type EngineHandler struct {
	engine    *engine.Engine
	stateRepo *postgres.StateRepo
}

func NewEngineHandler(e *engine.Engine, s *postgres.StateRepo) *EngineHandler {
	return &EngineHandler{
		engine:    e,
		stateRepo: s,
	}
}
func (h *EngineHandler) Next(w http.ResponseWriter, r *http.Request) {
	if h.engine == nil {
		panic("engine is nil")
	}

	if h.stateRepo == nil {
		panic("stateRepo is nil")
	}
	userID := r.URL.Query().Get("user")
	answer := r.URL.Query().Get("answer")

	current, _ := h.stateRepo.Get(userID)

	next := h.engine.Step(current, answer)

	if next == nil {
		w.Write([]byte("end"))
		return
	}

	h.stateRepo.Save(userID, next.ID)

	json.NewEncoder(w).Encode(next)
}
