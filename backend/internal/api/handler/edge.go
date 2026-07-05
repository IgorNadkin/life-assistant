package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domain/graph"
	"backend/internal/service"
)

type EdgeHandler struct {
	edgeService *service.EdgeService
}

func NewEdgeHandler(edgeService *service.EdgeService) *EdgeHandler {
	return &EdgeHandler{edgeService: edgeService}
}

type CreateEdgeRequest struct {
	FromNode  int64   `json:"from_node"`
	ToNode    int64   `json:"to_node"`
	Condition string  `json:"condition"`
	Logic     *string `json:"logic,omitempty"`
}

type EdgeResponse struct {
	ID        int64   `json:"id"`
	FromNode  int64   `json:"from_node"`
	ToNode    int64   `json:"to_node"`
	Condition string  `json:"condition"`
	Logic     *string `json:"logic,omitempty"`
}

// Create creates a new edge
func (h *EdgeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEdgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	edge := &graph.Edge{
		FromNode:  req.FromNode,
		ToNode:    req.ToNode,
		Condition: req.Condition,
		Logic:     req.Logic,
	}

	id, err := h.edgeService.CreateEdge(edge)
	if err != nil {
		http.Error(w, "failed to create edge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// Get retrieves an edge by ID
func (h *EdgeHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	edge, err := h.edgeService.GetEdge(id)
	if err != nil {
		http.Error(w, "edge not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(edge)
}

// GetByScenario retrieves all edges for a scenario
func (h *EdgeHandler) GetByScenario(w http.ResponseWriter, r *http.Request) {
	scenarioIDStr := r.URL.Query().Get("scenario_id")
	scenarioID, err := strconv.ParseInt(scenarioIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid scenario_id", http.StatusBadRequest)
		return
	}

	edges, err := h.edgeService.GetScenarioEdges(scenarioID)
	if err != nil {
		http.Error(w, "failed to retrieve edges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(edges)
}

// Update updates an edge
func (h *EdgeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req CreateEdgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	edge := &graph.Edge{
		ID:        id,
		FromNode:  req.FromNode,
		ToNode:    req.ToNode,
		Condition: req.Condition,
		Logic:     req.Logic,
	}

	if err := h.edgeService.UpdateEdge(edge); err != nil {
		http.Error(w, "failed to update edge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Delete deletes an edge
func (h *EdgeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.edgeService.DeleteEdge(id); err != nil {
		http.Error(w, "failed to delete edge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
