package service

import (
	"backend/internal/domain/scenario"
	"backend/internal/domain/user"
	"backend/internal/repository"
)

type ScenarioService struct {
	scenarioRepo repository.ScenarioRepository
	stateRepo    repository.UserStateRepository
	nodeRepo     repository.NodeRepository
}

func NewScenarioService(
	scenarioRepo repository.ScenarioRepository,
	stateRepo repository.UserStateRepository,
	nodeRepo repository.NodeRepository,
) *ScenarioService {
	return &ScenarioService{
		scenarioRepo: scenarioRepo,
		stateRepo:    stateRepo,
		nodeRepo:     nodeRepo,
	}
}

// CreateScenario creates a new scenario
func (s *ScenarioService) CreateScenario(scenario *scenario.Scenario) (int64, error) {
	return s.scenarioRepo.Create(scenario)
}

// GetScenario retrieves scenario by ID
func (s *ScenarioService) GetScenario(id int64) (*scenario.Scenario, error) {
	return s.scenarioRepo.Get(id)
}

// UpdateScenario updates an existing scenario
func (s *ScenarioService) UpdateScenario(scenario *scenario.Scenario) error {
	return s.scenarioRepo.Update(scenario)
}

// DeleteScenario deletes a scenario
func (s *ScenarioService) DeleteScenario(id int64) error {
	return s.scenarioRepo.Delete(id)
}

// ListScenarios retrieves all scenarios
func (s *ScenarioService) ListScenarios() ([]scenario.Scenario, error) {
	return s.scenarioRepo.List()
}

// StartScenario initializes a new user session for a scenario
func (s *ScenarioService) StartScenario(userID, scenarioID int64) (*user.UserState, error) {
	scen, err := s.scenarioRepo.Get(scenarioID)
	if err != nil {
		return nil, err
	}

	state := &user.UserState{
		UserID:         userID,
		ScenarioID:     scenarioID,
		CurrentNodeID:  scen.StartNodeID,
		Status:         user.StatusInProgress,
		CompletedSteps: user.Int64Array{},
	}

	id, err := s.stateRepo.Create(state)
	if err != nil {
		return nil, err
	}

	state.ID = id
	return state, nil
}

// GetUserScenarioProgress retrieves user's current progress in a scenario
func (s *ScenarioService) GetUserScenarioProgress(userID, scenarioID int64) (*user.UserState, error) {
	return s.stateRepo.GetByUserAndScenario(userID, scenarioID)
}
