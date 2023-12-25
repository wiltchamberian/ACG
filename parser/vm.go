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

	//global
	trueVal  NkObject
	falseVal NkObject
}

func NewVM(compiler NikaCompiler) *VM {
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

func (s *VM) popTwoIntegers() (int, int) {
	right := s.pop()
	left := s.pop()
	return left.(*NkInteger).Value, right.(*NkInteger).Value
}

func (s *VM) popInteger() int {
	return s.pop().(*NkInteger).Value
}

func (s *VM) pushInteger(val int) {
	s.push(&NkInteger{Value: val})
}

func (s *VM) pushBool(val bool) {
	if val {
		s.push(&NkInteger{Value: 1})
	} else {
		s.push(&NkInteger{Value: 0})
	}
}

func (s *VM) popBool() bool {
	return s.pop().(*NkInteger).Value != 0
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
				if err != nil {
					ip += InstructionByteLen(opCode)
					return obj, err
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
				s.push(s.trueVal)
			}
		case OpFalse:
			{
				s.push(s.falseVal)
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
		case OpBang:
			{
				l := s.popBool()
				s.pushBool(!l)
			}
		}
		ip += InstructionByteLen(opCode)
	}
	obj = s.StackTop()
	return obj, nil
}
