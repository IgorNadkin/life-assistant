package engine

import (
	"backend/internal/domain/graph"
	"backend/internal/domain/user"
)

type GraphEngine struct {
	nodes map[int64]graph.Node
	edges []graph.Edge
}

func New(nodes []graph.Node, edges []graph.Edge) *GraphEngine {
	m := make(map[int64]graph.Node)

	for _, n := range nodes {
		m[n.ID] = n
	}

	return &GraphEngine{
		nodes: m,
		edges: edges,
	}
}

// Step determines the next node based on current node and user answer
func (e *GraphEngine) Step(currentID int64, answer string) *graph.Node {
	for _, edge := range e.edges {
		if edge.FromNode != currentID {
			continue
		}

		if edge.Condition == answer {
			if node, ok := e.nodes[edge.ToNode]; ok {
				return &node
			}
		}
	}

	return nil
}

// GetCurrentNode returns the current node
func (e *GraphEngine) GetCurrentNode(nodeID int64) *graph.Node {
	if node, ok := e.nodes[nodeID]; ok {
		return &node
	}
	return nil
}

// UpdateState updates user state after a step
func (e *GraphEngine) UpdateState(state *user.UserState, nextNodeID int64) {
	state.CurrentNodeID = nextNodeID
	state.CompletedSteps = append(state.CompletedSteps, nextNodeID)
}

// IsTerminalNode checks if a node is a terminal (end) node
func (e *GraphEngine) IsTerminalNode(nodeID int64) bool {
	for _, edge := range e.edges {
		if edge.FromNode == nodeID {
			return false
		}
	}
	return true
}
