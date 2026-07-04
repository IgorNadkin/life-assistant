package repository

import "backend/internal/domain/scenario"

type ScenarioRepository interface {
	Create(s *scenario.Scenario) (int64, error)
	GetByID(id int64) (*scenario.Scenario, error)
}
