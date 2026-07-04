package graph

type Repository interface {
	GetNode(id int64) (Node, error)
	GetEdges(fromID int64) ([]Edge, error)
	UpdateSessionNode(sessionID int64, nodeID int64) error
	LogPath(sessionID int64, from, to int64) error
}
