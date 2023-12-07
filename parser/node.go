package parser

import (
	"bufio"
	"fmt"
	"os"
)

type INode interface {
	GetName() string
	GetChildren() []INode
}

type Node struct {
	name     string
	children []INode
}

func (s *Node) GetName() string {
	return s.name
}

func (s *Node) GetChildren() []INode {
	return s.children
}

func TreePrint(root INode) {

}

type Travel struct {
	level int
}

type NodeProcesser interface {
	ProcessNode(INode, int)
}

func (s *Travel) DepthFirstTravel(root INode, pc NodeProcesser) {
	pc.ProcessNode(root, s.level)
	s.level += 1
	children := root.GetChildren()
	for _, child := range children {
		s.DepthFirstTravel(child, pc)
	}
	s.level -= 1
}

type NodePrinter struct {
	file      *os.File
	writer    *bufio.Writer
	linebreak bool //flag if just print a "\n"
}

func (s *NodePrinter) Init(path string) (err error) {
	s.file, err = os.Create(path)
	if err != nil {
		return err
	}
	s.writer = bufio.NewWriter(s.file)
	return nil
}

func (s *NodePrinter) Close() {
	s.writer.Flush()
	s.file.Close()
}

func (s *NodePrinter) ProcessNode(root INode, level int) {
	if s.linebreak == true {
		for i := 0; i < level; i++ {
			fmt.Fprintf(s.writer, "             ")
		}
		s.linebreak = false
	}

	children := root.GetChildren()
	isLeaf := (children == nil) || (len(root.GetChildren()) == 0)
	if isLeaf {
		fmt.Fprintf(s.writer, "%10s\n", root.GetName())
		s.linebreak = true
	} else {
		fmt.Fprintf(s.writer, "%10s-->", root.GetName())
	}
}
