package parser

// a handwrite compilerl
type Compiler struct {
	instructions Instructions
	constants    []NkObject
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
