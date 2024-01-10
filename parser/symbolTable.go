package parser

type Scope = string

const (
	GlobalScope Scope = "GLOBAL"
)

// symbol
type Symbol struct {
	name  string
	scope Scope
	index int //variable index
}

type SymbolTable struct {
	mapping map[string]Symbol
	counter int
	globals *DataStack
}

func NewSymbolTable(globs *DataStack) *SymbolTable {
	return &SymbolTable{mapping: make(map[string]Symbol), counter: 0, globals: globs}
}

func (s *SymbolTable) def(node INode) Symbol {
	literal := node.GetLiteral()
	symbol := Symbol{literal, GlobalScope, s.globals.Length()}
	s.mapping[literal] = symbol
	s.globals.PushInteger(0)
	return symbol
}

func (s *SymbolTable) res(node INode) Symbol {
	return s.mapping[node.GetLiteral()]
}
