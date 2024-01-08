package parser

import (
	"encoding/binary"
	"fmt"
)

type OpCode = byte

type Instructions = []byte

type Definition struct {
	Name    string
	Widths  []int
	ByteLen int
}

const (
	OpConstant = iota
	OpAdd
	OpFAdd
	OpSub
	OpFSub
	OpMul
	OpDiv
	OpTrue
	OpFalse

	OpEq
	OpNotEq
	OpGt
	OpLe
	OpGtEq
	OpLeEq
	OpLShift
	OpRShift
	OpOr
	OpAnd

	OpBang //!
	OpBAnd //&
	OpBOr  //|
	OpPos  //+
	OpNeg  //-
	OpAddr //&

	OpGlobalSet
	OpGlobalGet

	OpJumpNotTrue //conditional jump
	OpJump        //jump

	OpPop  //pop the last element on stack
	OpPopd //pop the last element and put it on debug stack
)

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}, 3},
	OpAdd:      {"OpAdd", []int{}, 1},
	OpSub:      {"OpSub", []int{}, 1},
	OpMul:      {"OpMul", []int{}, 1},
	OpDiv:      {"OpDiv", []int{}, 1},
	OpTrue:     {"OpTrue", []int{}, 1},
	OpFalse:    {"OpTrue", []int{}, 1},
	OpEq:       {"OpEq", []int{}, 1},
	OpNotEq:    {"OpNotEq", []int{}, 1},
	OpLe:       {"OpLe", []int{}, 1},
	OpGt:       {"OpGt", []int{}, 1},
	OpLeEq:     {"OpLeEq", []int{}, 1},
	OpGtEq:     {"OpGtEq", []int{}, 1},
	OpLShift:   {"OpLShift", []int{}, 1},
	OpRShift:   {"OpRShift", []int{}, 1},
	OpOr:       {"OpOr", []int{}, 1},
	OpAnd:      {"OpAnd", []int{}, 1},
	OpBang:     {"OpBang", []int{}, 1},
	OpBAnd:     {"OpBAnd", []int{}, 1},
	OpBOr:      {"OpBOr", []int{}, 1},
	OpPos:      {"OpPos", []int{}, 1},
	OpNeg:      {"OpNeg", []int{}, 1},
	OpAddr:     {"OpAddr", []int{}, 1},

	OpGlobalSet: {"OpGlobalSet", []int{2}, 3},
	OpGlobalGet: {"OpGlobalGet", []int{2}, 3},

	OpJumpNotTrue: {"OpJumpNotTrue", []int{2}, 3},
	OpJump:        {"OpJump", []int{2}, 3},

	OpPop:  {"OpPop", []int{}, 1},
	OpPopd: {"OpPopd", []int{}, 1},
}

func LookUp(opcode OpCode) *Definition {
	val, err := definitions[opcode]
	if err != true {
		return nil
	}
	return val
}

func Make(opcode OpCode, operands ...int) []byte {
	definition, ok := definitions[opcode]
	if ok == false {
		return []byte{}
	}
	len := 1
	for _, variable := range definition.Widths {
		len += variable
	}
	bytes := make([]byte, len)
	bytes[0] = opcode
	offset := 1
	for i, operand := range operands {
		width := definition.Widths[i]
		switch width {
		case 2:
			binary.LittleEndian.PutUint16(bytes[offset:], uint16(operand))
		}
		offset += width
	}
	return bytes
}

func InstructionByteLen(opcode OpCode) int {
	definition := LookUp(opcode)
	if definition != nil {
		return definition.ByteLen
	}
	return 0
}

func M(nd INode) OpCode {
	switch nd.GetName() {
	case TkAdd:
		return OpAdd
	case TkSub:
		return OpSub
	case TkLShift:
		return OpLShift
	case TkRShift:
		return OpRShift
	case TkLess:
		return OpLe
	case TkGreater:
		return OpGt
	case TkLessEq:
		return OpLeEq
	case TkGreaterEq:
		return OpGtEq
	case TkMul:
		return OpMul
	case TkDiv:
		return OpDiv
	case TkEqual:
		return OpEq
	case TkNotEq:
		return OpNotEq
	case TkOr:
		return OpOr
	case TkAnd:
		return OpAnd
	}
	fmt.Println("OpCode fail")
	panic("")
}
