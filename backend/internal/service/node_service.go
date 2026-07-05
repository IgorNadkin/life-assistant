package service

import (
	"backend/internal/domain/graph"
	"backend/internal/repository"
)

type NodeService struct {
	nodeRepo repository.NodeRepository
}

func NewNodeService(nodeRepo repository.NodeRepository) *NodeService {
	return &NodeService{nodeRepo: nodeRepo}
}

// CreateNode creates a new graph node
func (s *NodeService) CreateNode(node *graph.Node) (int64, error) {
	return s.nodeRepo.Create(node)
}

// GetNode retrieves a node by ID
func (s *NodeService) GetNode(id int64) (*graph.Node, error) {
	return s.nodeRepo.GetByID(id)
}

// GetScenarioNodes retrieves all nodes for a scenario
func (s *NodeService) GetScenarioNodes(scenarioID int64) ([]graph.Node, error) {
	return s.nodeRepo.GetByScenario(scenarioID)
}

// UpdateNode updates an existing node
func (s *NodeService) UpdateNode(node *graph.Node) error {
	return s.nodeRepo.Update(node)
}

// DeleteNode deletes a node
func (s *NodeService) DeleteNode(id int64) error {
	return s.nodeRepo.Delete(id)
}
