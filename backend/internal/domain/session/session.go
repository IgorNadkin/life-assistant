package session

type Session struct {
	ID            int64
	UserID        int64
	ScenarioID    int64
	CurrentNodeID int64
	Status        string
}
