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

	//global variable
	globals []NkObject

	//global
	trueVal  NkObject
	falseVal NkObject
}

func NewVM(compiler *NikaCompiler) *VM {
	var vm VM
	vm.constants = compiler.constants
	vm.instructions = compiler.instructions
	vm.sp = 0
	vm.stack = make([]NkObject, 100)

	vm.trueVal = &NkInteger{Value: 1}
	vm.falseVal = &NkInteger{Value: 0}
	return &vm
}

func (s *VM) StackTop() NkObject {
	if s.sp == 0 {
		return nil
	}
	return s.stack[s.sp-1]
}

func (s *VM) Pop() NkObject {
	if s.sp > 0 {
		s.sp--
		return s.stack[s.sp]
	}
	return nil
}

func (s *VM) Push(o NkObject) error {
	if s.sp >= len(s.stack) {
		return errors.New("stack overflow")
	}
	s.stack[s.sp] = o
	s.sp++
	return nil
}

func (s *VM) popTwoIntegers() (int, int) {
	right := s.Pop()
	left := s.Pop()
	return left.(*NkInteger).Value, right.(*NkInteger).Value
}

func (s *VM) popInteger() int {
	return s.Pop().(*NkInteger).Value
}

func (s *VM) pushInteger(val int) {
	s.Push(&NkInteger{Value: val})
}

func (s *VM) pushBool(val bool) {
	if val {
		s.Push(&NkInteger{Value: 1})
	} else {
		s.Push(&NkInteger{Value: 0})
	}
}

func (s *VM) popBool() bool {
	return s.Pop().(*NkInteger).Value != 0
}

func (s *VM) IsEmpty() bool {
	return s.sp <= 0
}

func (s *VM) Run() error {
	var err error
	ip := 0
	l := len(s.instructions)
	for ip < l {
		opCode := s.instructions[ip]
		switch opCode {
		case OpConstant:
			{
				index := ReadUint16(s.instructions[ip+1:])
				err = s.Push(s.constants[index])
				if err != nil {
					ip += InstructionByteLen(opCode)
					return err
				}
			}
		case OpAdd:
			{
				l, r := s.popTwoIntegers()
				s.pushInteger(l + r)
			}
		case OpSub:
			{
				l, r := s.popTwoIntegers()
				s.pushInteger(l - r)
			}
		case OpMul:
			{
				l, r := s.popTwoIntegers()
				s.pushInteger(l * r)
			}
		case OpDiv:
			{
				l, r := s.popTwoIntegers()
				s.pushInteger(l / r)
			}
		case OpTrue:
			{
				s.Push(s.trueVal)
			}
		case OpFalse:
			{
				s.Push(s.falseVal)
			}
		case OpEq:
			{
				l, r := s.popTwoIntegers()
				s.pushBool(l == r)
			}
		case OpNotEq:
			{
				l, r := s.popTwoIntegers()
				s.pushBool(l != r)
			}
		case OpLe:
			{
				l, r := s.popTwoIntegers()
				s.pushBool(l < r)
			}
		case OpGt:
			{
				l, r := s.popTwoIntegers()
				s.pushBool(l > r)
			}
		case OpLeEq:
			{
				l, r := s.popTwoIntegers()
				s.pushBool(l <= r)
			}
		case OpGtEq:
			{
				l, r := s.popTwoIntegers()
				s.pushBool(l >= r)
			}
		case OpNeg:
			{
				l := s.popInteger()
				s.pushInteger(-l)
			}
		case OpBang:
			{
				l := s.popBool()
				s.pushBool(!l)
			}
		case OpOr:
			{
				l, r := s.popTwoIntegers()
				s.pushBool((l != 0) || (r != 0))
			}
		case OpAnd:
			{
				l, r := s.popTwoIntegers()
				s.pushBool((l != 0) && (r != 0))
			}
		case OpGlobalSet:
			{
				index := ReadUint16(s.instructions[ip+1:])
				s.globals[index] = s.Pop()
			}
		case OpGlobalGet:
			{
				index := ReadUint16(s.instructions[ip+1:])
				s.Push(s.globals[index])
			}
		}
		ip += InstructionByteLen(opCode)
	}
	return nil
}
