package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domain/scenario"
	"backend/internal/repository/postgres"
)

type ScenarioHandler struct {
	repo *postgres.ScenarioRepo
}

func NewScenarioHandler(repo *postgres.ScenarioRepo) *ScenarioHandler {
	return &ScenarioHandler{repo: repo}
}

func (h *ScenarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s scenario.Scenario
	json.NewDecoder(r.Body).Decode(&s)

	id, _ := h.repo.Create(&s)

	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *ScenarioHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	s, _ := h.repo.Get(id)

	json.NewEncoder(w).Encode(s)
}
