package parser

import "strconv"

//a handwrite compilerl
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

func (s *Compiler) Compile(node INode) error {
	name := node.GetName()
	switch name {
	case "Number":
		{
			d, er := strconv.Atoi(node.GetLiteral())
			if er != nil {
				return er
			}
			val := &NkInteger{Value: d}
			s.emit(OpConstant, s.addConstant(val))
		}
	case "expr":
		{
			var err error
			v := node.GetChildren()
			if nil == s.Compile(v[0]) &&
				nil == s.Compile(v[2]) {
			} else {
				return err
			}

			switch v[1].GetName() {
			case "+":
				{
					s.emit(OpAdd)
				}
			}
		}
	case "stmt":
		{
			child := node.GetChildren()[0]
			if child.GetName() == "expr" {

			}
		}
	}

	return nil
}
