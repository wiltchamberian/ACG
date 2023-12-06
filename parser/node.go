package parser

type INode interface {
	GetChildren() []INode
}

type Node struct {
	name     string
	children []INode
}

func (s *Node) GetChildren() []INode {
	return s.children
}
