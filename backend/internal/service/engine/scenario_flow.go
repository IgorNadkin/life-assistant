package engine

import (
	"backend/internal/domain/graph"
	"backend/internal/domain/user"
	"backend/internal/repository"
)

type ScenarioFlow struct {
	engine     *GraphEngine
	stateRepo  repository.UserStateRepository
	actionRepo repository.UserActionRepository
}

func NewScenarioFlow(
	engine *GraphEngine,
	stateRepo repository.UserStateRepository,
	actionRepo repository.UserActionRepository,
) *ScenarioFlow {
	return &ScenarioFlow{
		engine:     engine,
		stateRepo:  stateRepo,
		actionRepo: actionRepo,
	}
}

// ProcessAnswer processes user's answer and moves to next node
func (sf *ScenarioFlow) ProcessAnswer(state *user.UserState, answer string) (*graph.Node, error) {
	nextNode := sf.engine.Step(state.CurrentNodeID, answer)
	if nextNode == nil {
		// Terminal node reached
		state.Status = user.StatusCompleted
		if err := sf.stateRepo.Update(state); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Update state
	sf.engine.UpdateState(state, nextNode.ID)
	if err := sf.stateRepo.Update(state); err != nil {
		return nil, err
	}

	return nextNode, nil
}

// GetCurrentNode returns the current node for user state
func (sf *ScenarioFlow) GetCurrentNode(state *user.UserState) *graph.Node {
	return sf.engine.GetCurrentNode(state.CurrentNodeID)
}

// CreateActionFromNode creates a user action from node information (if node type is 'action')
func (sf *ScenarioFlow) CreateActionFromNode(stateID int64, node *graph.Node) (*user.UserAction, error) {
	if node.Type != graph.NodeTypeAction {
		return nil, nil // Not an action node
	}

	action := &user.UserAction{
		UserStateID:  stateID,
		NodeID:       node.ID,
		Action:       node.ActionType,
		Organization: *node.Organization,
		Deadline:     nil, // Could be parsed from node.Deadline
		Completed:    false,
	}

	id, err := sf.actionRepo.Create(action)
	if err != nil {
		return nil, err
	}

	action.ID = id
	return action, nil
}
