package engine

import (
	"backend/internal/domain/graph"
)

type Engine struct {
	nodes map[int64]graph.Node
	edges []graph.Edge
}

func New(nodes []graph.Node, edges []graph.Edge) *Engine {
	m := make(map[int64]graph.Node)

	for _, n := range nodes {
		m[n.ID] = n
	}

	return &Engine{
		nodes: m,
		edges: edges,
	}
}
func (e *Engine) Step(currentID int64, answer string) *graph.Node {
	for _, edge := range e.edges {
		if edge.From != currentID {
			continue
		}

		if edge.Condition == answer {
			node := e.nodes[edge.To]
			return &node
		}
	}

	return nil
}
