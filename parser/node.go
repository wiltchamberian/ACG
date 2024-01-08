package parser

import (
	"fmt"
	"slices"
	"strconv"
)

type INode interface {
	SetName(string)
	GetName() string
	SetType(tp NikaType)
	GetType() *NikaType
	GetInteger() int32
	GetLiteral() string
	GetChildren() []INode
	GetParent() INode
	AddChild(INode)
	SetParent(INode)
	IsTerminal() bool
	Select() int //记录该节点在展开时走到哪个alternative分支
	//Action() string //记录该节点展开时走到那个分支的action

	Compile(s *NikaCompiler) error
}

func NodeCount(inode INode) int {
	var count int = 1
	children := inode.GetChildren()
	for _, v := range children {
		count += NodeCount(v)
	}
	return count
}

func IsLeafNode(nd INode) bool {
	return len(nd.GetChildren()) == 0
}

func ToNkInteger(nd INode) *NkInteger {
	d, er := strconv.Atoi(nd.GetLiteral())
	if er != nil {
		return nil
	}
	return &NkInteger{Value: d}
}

func ToInteger(nd INode) int32 {
	d, er := strconv.Atoi(nd.GetLiteral())
	if er != nil {
		return 0
	}
	return int32(d)
}

type Node struct {
	Name     string
	Children []INode
	Parent   INode
	selected int
	Action   string
	NkType   NikaType
}

func (s *Node) SetType(tp NikaType) {
	s.NkType = tp
}

func (s *Node) GetType() *NikaType {
	return &s.NkType
}

func (s *Node) GetLiteral() string {
	return ""
}

func (s *Node) SetName(name string) {
	s.Name = name
}

func (s *Node) GetName() string {
	return s.Name
}

func (s *Node) GetChildren() []INode {
	return s.Children
}

func (s *Node) GetParent() INode {
	return s.Parent
}

func (s *Node) AddChild(child INode) {
	s.Children = append(s.Children, child)
}

func (s *Node) SetParent(parent INode) {
	s.Parent = parent
}

func (s *Node) Select() int {
	return s.selected
}

func (s *Node) IsTerminal() bool {
	return false
}

func (s *Node) Compile(c *NikaCompiler) error {
	return nil
}

func (s *Node) GetInteger() int32 {
	d, _ := strconv.Atoi(s.GetLiteral())
	return int32(d)
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
	FileWriter
	linebreak bool //flag if just print a "\n"
}

func (s *NodePrinter) ProcessNode(root INode, level int) {
	if s.linebreak == true {
		for i := 0; i < level; i++ {
			fmt.Fprintf(s.writer, "         ")
		}
		s.linebreak = false
	}

	children := root.GetChildren()
	isLeaf := (children == nil) || (len(root.GetChildren()) == 0)
	if isLeaf {
		fmt.Fprintf(s.writer, "%6s\n", root.GetLiteral())
		s.linebreak = true
	} else {
		fmt.Fprintf(s.writer, "%6s-->", root.GetName())
	}
}

func ToAST(nd INode) {
	name := nd.GetName()
	children := nd.GetChildren()
	if name == "term" || name == "expr" || name == "assign" {
		n := nd.(*Node)
		if len(children) == 3 {
			n.Name = n.Children[1].GetName()
			n.Children = slices.Delete(n.Children, 1, 2)
		} else if len(children) == 1 {
			n.Name = n.Children[0].GetName()
			n.Children = n.Children[0].GetChildren()
			ToAST(n)
			return
		} else {
			panic("")
		}
	} else if name == "atom" {
		n := nd.(*Node)
		if len(children) == 3 {
			n.Name = n.Children[1].GetName()
			n.Children = n.Children[1].GetChildren()
			ToAST(n)
			return
		} else if len(children) == 1 {
			n.Name = n.Children[0].GetName()
			n.Children = n.Children[0].GetChildren()
			ToAST(n)
			return
		} else {
			panic("")
		}
	} else if name == "if" {
		n := nd.(*Node)
		if len(children) >= 5 {
			arr := []INode{n.Children[1], n.Children[3]}
			n.Children = arr
		} else {
			panic("")
		}
	} else if name == "for" {
		n := nd.(*Node)
		if len(children) >= 9 {
			arr := []INode{n.Children[1], n.Children[3], n.Children[5], n.Children[7]}
			n.Children = arr
		} else {
			panic("")
		}
	}
	children = nd.GetChildren() //update
	for _, child := range children {
		ToAST(child)
	}

}

type Ret struct {
	Nd  INode
	Err error
}

// helper
func Append(nodes *[]INode, ret Ret) bool {
	if ret.Err == nil {
		*nodes = append(*nodes, ret.Nd)
		return true
	}
	return false
}

func NewEmpty() Ret {
	return Ret{Nd: EmptyToken, Err: nil}
}
