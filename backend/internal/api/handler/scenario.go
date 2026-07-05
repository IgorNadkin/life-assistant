package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domain/scenario"
	"backend/internal/service"
)

type ScenarioHandler struct {
	scenarioService *service.ScenarioService
}

func NewScenarioHandler(scenarioService *service.ScenarioService) *ScenarioHandler {
	return &ScenarioHandler{scenarioService: scenarioService}
}

type CreateScenarioRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    *string `json:"category,omitempty"`
	StartNodeID int64   `json:"start_node_id"`
}

type ScenarioResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    *string `json:"category,omitempty"`
	StartNodeID int64   `json:"start_node_id"`
}

// Create creates a new scenario
func (h *ScenarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateScenarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	s := &scenario.Scenario{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		StartNodeID: req.StartNodeID,
	}

	id, err := h.scenarioService.CreateScenario(s)
	if err != nil {
		http.Error(w, "failed to create scenario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// Get retrieves a scenario by ID
func (h *ScenarioHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	s, err := h.scenarioService.GetScenario(id)
	if err != nil {
		http.Error(w, "scenario not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// Update updates a scenario
func (h *ScenarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req CreateScenarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	s := &scenario.Scenario{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		StartNodeID: req.StartNodeID,
	}

	if err := h.scenarioService.UpdateScenario(s); err != nil {
		http.Error(w, "failed to update scenario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Delete deletes a scenario
func (h *ScenarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.scenarioService.DeleteScenario(id); err != nil {
		http.Error(w, "failed to delete scenario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// List retrieves all scenarios
func (h *ScenarioHandler) List(w http.ResponseWriter, r *http.Request) {
	scenarios, err := h.scenarioService.ListScenarios()
	if err != nil {
		http.Error(w, "failed to list scenarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenarios)
}
