package parser

import "fmt"

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
	literal := node.GetLiteral()
	symbol := Symbol{literal, GlobalScope, s.counter}
	s.mapping[literal] = symbol
	s.counter++
	return symbol
}

func (s *SymbolTable) res(node INode) Symbol {
	return s.mapping[node.GetLiteral()]
}

// a handwrite compilerl
type Compiler struct {
	instructions Instructions
	constants    []NkObject
	symbolTable  *SymbolTable

	poses []int
}

//push index of instruction
func (s *Compiler) pushIndex(index int) {
	s.poses = append(s.poses, index)
}

func (s *Compiler) popIndex() int {
	x := s.poses[len(s.poses)-1]
	s.poses = s.poses[0 : len(s.poses)-1]
	return x
}

func (s *Compiler) InitCompiler() {
	s.symbolTable = NewSymbolTable()
	s.poses = make([]int, 0)
}

func (s *Compiler) addConstant(obj NkObject) int {
	pos := len(s.constants)
	s.constants = append(s.constants, obj)
	return pos
}

func (s *Compiler) emit(opcode OpCode, operands ...int) int {
	//for test
	fmt.Printf("emit:%d\n", opcode)
	ins := Make(opcode, operands...)
	pos := s.addInstruction(ins)
	return pos
}

func (s *Compiler) pop() {
	s.emit(OpPop)
}

func (s *Compiler) popd() {
	s.emit(OpPopd)
}

func (s *Compiler) Pos() int {
	return len(s.instructions)
}

//jump not true
func (s *Compiler) jumpNt() int {
	return s.emit(OpJumpNotTrue, s.currentPos())
}

func (s *Compiler) jump() int {
	return s.emit(OpJump, s.currentPos())
}

func (s *Compiler) jumpTo(pos int) int {
	return s.emit(OpJump, pos)
}

func (s *Compiler) addInstruction(ins []byte) int {
	pos := len(s.instructions)
	s.instructions = append(s.instructions, ins...)
	return pos
}

func (s *Compiler) def(node INode) Symbol {
	return s.symbolTable.def(node)
}

func (s *Compiler) res(node INode) Symbol {
	return s.symbolTable.res(node)
}

func (s *Compiler) replaceInstruction(pos int, ins []byte) {
	for i := 0; i < len(ins); i++ {
		s.instructions[pos+i] = ins[i]
	}
}

func (s *Compiler) replace(pos int) {
	op := s.instructions[pos]
	newIns := Make(op, len(s.instructions))
	s.replaceInstruction(pos, newIns)
}

func (s *Compiler) replaceAll() {
	l := len(s.poses)
	for i := 0; i < l; i++ {
		s.replace(s.popIndex())
	}
}

func (s *Compiler) currentPos() int {
	return len(s.instructions)
}
