package graph

type NodeType string

const (
	NodeQuestion NodeType = "question"
	NodeAction   NodeType = "action"
	NodeEnd      NodeType = "end"
)

type Node struct {
	ID   int64
	Type NodeType
	Text string
}
