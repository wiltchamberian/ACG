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
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{mapping: make(map[string]Symbol), counter: 0}
}

func (s *SymbolTable) def(node INode) Symbol {
	symbol := Symbol{node.GetName(), GlobalScope, s.counter}
	s.mapping[node.GetName()] = symbol
	s.counter++
	return symbol
}

func (s *SymbolTable) res(node INode) Symbol {
	return s.mapping[node.GetName()]
}

// a handwrite compilerl
type Compiler struct {
	instructions Instructions
	constants    []NkObject
	symbolTable  *SymbolTable
}

func (compiler *Compiler) InitCompiler() {
	compiler.symbolTable = NewSymbolTable()
}

func (s *Compiler) addConstant(obj NkObject) int {
	pos := len(s.constants)
	s.constants = append(s.constants, obj)
	return pos
}

func (s *Compiler) emit(opcode OpCode, operands ...int) int {
	ins := Make(opcode, operands...)
	pos := s.addInstruction(ins)
	return pos
}

func (s *Compiler) addInstruction(ins []byte) int {
	pos := len(s.instructions)
	s.instructions = append(s.instructions, ins...)
	return pos
}

func (s *Compiler) def(node INode) {
	s.symbolTable.def(node)
}

func (s *Compiler) res(node INode) {
	s.symbolTable.res(node)
}
