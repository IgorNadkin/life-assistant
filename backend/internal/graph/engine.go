package graph

type Engine struct {
	repo Repository
}

func NewEngine(repo Repository) *Engine {
	return &Engine{repo: repo}
}

func (e *Engine) Next(sessionID int64, answers map[string]any, current Node) (*Node, error) {

	edges, err := e.repo.GetEdges(current.ID)
	if err != nil {
		return nil, err
	}

	for _, edge := range edges {

		if Evaluate(edge.Condition, answers) {

			e.repo.UpdateSessionNode(sessionID, edge.ToNodeID)
			e.repo.LogPath(sessionID, current.ID, edge.ToNodeID)

			node, err := e.repo.GetNode(edge.ToNodeID)
			if err != nil {
				return nil, err
			}

			return &node, nil
		}
	}

	return nil, nil
}
