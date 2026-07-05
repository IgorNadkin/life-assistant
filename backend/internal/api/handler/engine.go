package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domain/user"
	"backend/internal/service"
	"backend/internal/service/engine"
)

type EngineHandler struct {
	scenarioService *service.ScenarioService
	userService     *service.UserService
	scenarioFlow    *engine.ScenarioFlow
}

func NewEngineHandler(
	scenarioService *service.ScenarioService,
	userService *service.UserService,
	scenarioFlow *engine.ScenarioFlow,
) *EngineHandler {
	return &EngineHandler{
		scenarioService: scenarioService,
		userService:     userService,
		scenarioFlow:    scenarioFlow,
	}
}

type StartScenarioRequest struct {
	UserID     int64 `json:"user_id"`
	ScenarioID int64 `json:"scenario_id"`
}

type StepRequest struct {
	UserStateID int64  `json:"user_state_id"`
	Answer      string `json:"answer"`
}

type NodeResponse struct {
	ID             int64   `json:"id"`
	Type           string  `json:"type"`
	Text           string  `json:"text"`
	ActionType     *string `json:"action_type,omitempty"`
	Organization   *string `json:"organization,omitempty"`
	Deadline       *string `json:"deadline,omitempty"`
	ReferenceLinks *string `json:"reference_links,omitempty"`
	Consequences   *string `json:"consequences,omitempty"`
}

type StepResponse struct {
	CurrentNode *NodeResponse `json:"current_node,omitempty"`
	IsCompleted bool          `json:"is_completed"`
	Message     string        `json:"message,omitempty"`
}

// Start initializes a new user session for a scenario
func (h *EngineHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req StartScenarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	state, err := h.scenarioService.StartScenario(req.UserID, req.ScenarioID)
	if err != nil {
		http.Error(w, "failed to start scenario", http.StatusInternalServerError)
		return
	}

	currentNode := h.scenarioFlow.GetCurrentNode(state)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"state_id":      state.ID,
		"current_node":  currentNode,
		"status":        state.Status,
	})
}

// Step processes user answer and returns next node
func (h *EngineHandler) Step(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	state, err := h.scenarioFlow.(*engine.ScenarioFlow).stateRepo.GetByID(req.UserStateID)
	if err != nil {
		http.Error(w, "state not found", http.StatusNotFound)
		return
	}

	// This is a workaround - in production you'd expose stateRepo properly
	// For now, we'll use the service to get state
	if state == nil {
		http.Error(w, "state not found", http.StatusNotFound)
		return
	}

	nextNode, err := h.scenarioFlow.ProcessAnswer(state, req.Answer)
	if err != nil {
		http.Error(w, "failed to process answer", http.StatusInternalServerError)
		return
	}

	response := StepResponse{
		IsCompleted: state.Status == user.StatusCompleted,
	}

	if nextNode != nil {
		response.CurrentNode = &NodeResponse{
			ID:             nextNode.ID,
			Type:           string(nextNode.Type),
			Text:           nextNode.Text,
			ActionType:     nextNode.ActionType,
			Organization:   nextNode.Organization,
			Deadline:       nextNode.Deadline,
			ReferenceLinks: nextNode.ReferenceLinks,
			Consequences:   nextNode.Consequences,
		}

		// If this is an action node, create an action record
		if nextNode.Type == "action" {
			_, _ = h.scenarioFlow.CreateActionFromNode(state.ID, nextNode)
		}
	} else if response.IsCompleted {
		response.Message = "Scenario completed!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetStatus retrieves current state of user in scenario
func (h *EngineHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	scenarioIDStr := r.URL.Query().Get("scenario_id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	scenarioID, err := strconv.ParseInt(scenarioIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid scenario_id", http.StatusBadRequest)
		return
	}

	state, err := h.scenarioService.GetUserScenarioProgress(userID, scenarioID)
	if err != nil {
		http.Error(w, "state not found", http.StatusNotFound)
		return
	}

	actions, _ := h.userService.GetUserActions(state.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"state":   state,
		"actions": actions,
	})
}
