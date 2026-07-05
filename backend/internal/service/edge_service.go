package service

import (
	"backend/internal/domain/graph"
	"backend/internal/repository"
)

type EdgeService struct {
	edgeRepo repository.EdgeRepository
}

func NewEdgeService(edgeRepo repository.EdgeRepository) *EdgeService {
	return &EdgeService{edgeRepo: edgeRepo}
}

// CreateEdge creates a new graph edge
func (s *EdgeService) CreateEdge(edge *graph.Edge) (int64, error) {
	return s.edgeRepo.Create(edge)
}

// GetEdge retrieves an edge by ID
func (s *EdgeService) GetEdge(id int64) (*graph.Edge, error) {
	return s.edgeRepo.GetByID(id)
}

// GetScenarioEdges retrieves all edges for a scenario
func (s *EdgeService) GetScenarioEdges(scenarioID int64) ([]graph.Edge, error) {
	return s.edgeRepo.GetByScenario(scenarioID)
}

// UpdateEdge updates an existing edge
func (s *EdgeService) UpdateEdge(edge *graph.Edge) error {
	return s.edgeRepo.Update(edge)
}

// DeleteEdge deletes an edge
func (s *EdgeService) DeleteEdge(id int64) error {
	return s.edgeRepo.Delete(id)
}
