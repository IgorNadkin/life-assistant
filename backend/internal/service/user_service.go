package service

import (
	"backend/internal/domain/user"
	"backend/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	stateRepo repository.UserStateRepository
	actionRepo repository.UserActionRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	stateRepo repository.UserStateRepository,
	actionRepo repository.UserActionRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		stateRepo: stateRepo,
		actionRepo: actionRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(u *user.User) (int64, error) {
	return s.userRepo.Create(u)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int64) (*user.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*user.User, error) {
	return s.userRepo.GetByEmail(email)
}

// GetUserScenarioState retrieves user's state for a specific scenario
func (s *UserService) GetUserScenarioState(userID, scenarioID int64) (*user.UserState, error) {
	return s.stateRepo.GetByUserAndScenario(userID, scenarioID)
}

// GetUserActions retrieves all actions for a user state
func (s *UserService) GetUserActions(userStateID int64) ([]user.UserAction, error) {
	return s.actionRepo.GetByUserState(userStateID)
}

// CompleteAction marks an action as completed
func (s *UserService) CompleteAction(action *user.UserAction) error {
	return s.actionRepo.Update(action)
}
