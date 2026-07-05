package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"backend/internal/domain/graph"
	"backend/internal/service"
)

type NodeHandler struct {
	nodeService *service.NodeService
}

func NewNodeHandler(nodeService *service.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

type CreateNodeRequest struct {
	ScenarioID     int64   `json:"scenario_id"`
	Type           string  `json:"type"`
	Text           string  `json:"text"`
	ActionType     *string `json:"action_type,omitempty"`
	Organization   *string `json:"organization,omitempty"`
	Deadline       *string `json:"deadline,omitempty"`
	ReferenceLinks *string `json:"reference_links,omitempty"`
	Consequences   *string `json:"consequences,omitempty"`
	Order          int     `json:"order"`
}

// backend/internal/api/handler/node.go
type NodeDetailResponse struct { // было NodeResponse
	ID             int64   `json:"id"`
	ScenarioID     int64   `json:"scenario_id"`
	Type           string  `json:"type"`
	Text           string  `json:"text"`
	ActionType     *string `json:"action_type,omitempty"`
	Organization   *string `json:"organization,omitempty"`
	Deadline       *string `json:"deadline,omitempty"`
	ReferenceLinks *string `json:"reference_links,omitempty"`
	Consequences   *string `json:"consequences,omitempty"`
	Order          int     `json:"order"`
}

// Create creates a new node
func (h *NodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	node := &graph.Node{
		ScenarioID:     req.ScenarioID,
		Type:           graph.NodeType(req.Type),
		Text:           req.Text,
		ActionType:     req.ActionType,
		Organization:   req.Organization,
		Deadline:       req.Deadline,
		ReferenceLinks: req.ReferenceLinks,
		Consequences:   req.Consequences,
		Order:          req.Order,
	}

	id, err := h.nodeService.CreateNode(node)
	if err != nil {
		log.Printf("[NodeHandler.Create] %v", err)
		http.Error(w, "failed to create node", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// Get retrieves a node by ID
func (h *NodeHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	node, err := h.nodeService.GetNode(id)
	if err != nil {
		http.Error(w, "node not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

// GetByScenario retrieves all nodes for a scenario
func (h *NodeHandler) GetByScenario(w http.ResponseWriter, r *http.Request) {
	scenarioIDStr := r.URL.Query().Get("scenario_id")
	scenarioID, err := strconv.ParseInt(scenarioIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid scenario_id", http.StatusBadRequest)
		return
	}

	nodes, err := h.nodeService.GetScenarioNodes(scenarioID)
	if err != nil {
		http.Error(w, "failed to retrieve nodes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

// Update updates a node
func (h *NodeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	node := &graph.Node{
		ID:             id,
		ScenarioID:     req.ScenarioID,
		Type:           graph.NodeType(req.Type),
		Text:           req.Text,
		ActionType:     req.ActionType,
		Organization:   req.Organization,
		Deadline:       req.Deadline,
		ReferenceLinks: req.ReferenceLinks,
		Consequences:   req.Consequences,
		Order:          req.Order,
	}

	if err := h.nodeService.UpdateNode(node); err != nil {
		http.Error(w, "failed to update node", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Delete deletes a node
func (h *NodeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.nodeService.DeleteNode(id); err != nil {
		http.Error(w, "failed to delete node", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
