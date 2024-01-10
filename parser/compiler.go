package parser

import "fmt"

// a handwrite compilerl
type Compiler struct {
	instructions    Instructions
	constantObjects []NkObject
	constants       *DataStack
	globals         *DataStack
	symbolTable     *SymbolTable
	typeSystem      *TypeSystem
	poses           []int
}

func (s *Compiler) NewStruct(name string, fields []FieldShow) bool {
	return s.typeSystem.NewStruct(name, fields)
}

func (s *Compiler) NewInteger() bool {
	return s.typeSystem.NewInteger()
}

func (s *Compiler) NewFloat() bool {
	return s.typeSystem.NewFloat()
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
	s.globals = NewDataStack()
	s.constants = NewDataStack()
	s.symbolTable = NewSymbolTable(s.globals)
	s.poses = make([]int, 0)
}

func (s *Compiler) addConstantObject(obj NkObject) int {
	pos := len(s.constantObjects)
	s.constantObjects = append(s.constantObjects, obj)
	return pos
}

func (s *Compiler) addConstantInteger(val int32) int {
	return s.constants.PushInteger(val)
}

func (s *Compiler) emit(opcode OpCode, operands ...int) int {
	//for test
	fmt.Printf("emit:%d\n", opcode)
	ins := Make(opcode, operands...)
	pos := s.addInstruction(ins)
	return pos
}

func (s *Compiler) pop(node INode) {
	s.emit(OpPop)
}

func (s *Compiler) popd(node INode) {
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
