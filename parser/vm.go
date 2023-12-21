package parser

import (
	"encoding/binary"
	"errors"
)

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

type VM struct {
	constants    []NkObject
	instructions Instructions
	stack        []NkObject
	sp           int
}

func NewVM(compiler NikaCompiler) *VM {
	var vm VM
	vm.constants = compiler.constants
	vm.instructions = compiler.instructions
	vm.sp = 0
	vm.stack = make([]NkObject, 100)
	return &vm
}

func (s *VM) StackTop() NkObject {
	if s.sp == 0 {
		return nil
	}
	return s.stack[s.sp-1]
}

func (s *VM) pop() NkObject {
	o := s.stack[s.sp-1]
	s.sp--
	return o
}

func (s *VM) push(o NkObject) error {
	if s.sp >= len(s.stack) {
		return errors.New("stack overflow")
	}
	s.stack[s.sp] = o
	s.sp++
	return nil
}

func (s *VM) Run() (NkObject, error) {
	var err error
	var obj NkObject
	ip := 0
	l := len(s.instructions)
	for ip < l {
		opCode := s.instructions[ip]
		switch opCode {
		case OpConstant:
			{
				index := ReadUint16(s.instructions[ip+1:])
				err = s.push(s.constants[index])
				ip += InstructionByteLen(opCode)
				if err != nil {
					return obj, err
				}
			}
		case OpAdd:
			{
				right := s.pop()
				left := s.pop()
				rightValue := right.(*NkInteger).Value
				leftValue := left.(*NkInteger).Value
				result := leftValue + rightValue
				s.push(&NkInteger{Value: result})
				ip += InstructionByteLen(opCode)
			}
		case OpSub:
			{
				right := s.pop()
				left := s.pop()
				rightValue := right.(*NkInteger).Value
				leftValue := left.(*NkInteger).Value
				result := leftValue - rightValue
				s.push(&NkInteger{Value: result})
				ip += InstructionByteLen(opCode)
			}
		case OpMul:
			{
				right := s.pop()
				left := s.pop()
				rightValue := right.(*NkInteger).Value
				leftValue := left.(*NkInteger).Value
				result := leftValue * rightValue
				s.push(&NkInteger{Value: result})
				ip += InstructionByteLen(opCode)
			}
		case OpDiv:
			{
				right := s.pop()
				left := s.pop()
				rightValue := right.(*NkInteger).Value
				leftValue := left.(*NkInteger).Value
				result := leftValue / rightValue
				s.push(&NkInteger{Value: result})
				ip += InstructionByteLen(opCode)
			}
		}
	}
	obj = s.StackTop()
	return obj, nil
}
