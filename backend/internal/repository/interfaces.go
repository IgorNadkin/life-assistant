package repository

import (
	"backend/internal/domain/graph"
	"backend/internal/domain/scenario"
	"backend/internal/domain/user"
)

type ScenarioRepository interface {
	Create(s *scenario.Scenario) (int64, error)
	Get(id int64) (*scenario.Scenario, error)
	Update(s *scenario.Scenario) error
	Delete(id int64) error
	List() ([]scenario.Scenario, error)
}

type NodeRepository interface {
	Create(node *graph.Node) (int64, error)
	GetByID(id int64) (*graph.Node, error)
	GetByScenario(scenarioID int64) ([]graph.Node, error)
	Update(node *graph.Node) error
	Delete(id int64) error
}

type EdgeRepository interface {
	Create(edge *graph.Edge) (int64, error)
	GetByID(id int64) (*graph.Edge, error)
	GetByScenario(scenarioID int64) ([]graph.Edge, error)
	Update(edge *graph.Edge) error
	Delete(id int64) error
}

type UserRepository interface {
	Create(u *user.User) (int64, error)
	GetByID(id int64) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
}

type UserStateRepository interface {
	Create(state *user.UserState) (int64, error)
	GetByID(id int64) (*user.UserState, error)
	GetByUser(userID int64) (*user.UserState, error)
	GetByUserAndScenario(userID, scenarioID int64) (*user.UserState, error)
	Update(state *user.UserState) error
}

type UserActionRepository interface {
	Create(action *user.UserAction) (int64, error)
	GetByUserState(userStateID int64) ([]user.UserAction, error)
	Update(action *user.UserAction) error
}
