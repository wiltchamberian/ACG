package parser

import (
	"encoding/binary"
)

const (
	OpConstant = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
)

type OpCode = byte

type Instructions = []byte

type Definition struct {
	Name    string
	Widths  []int
	ByteLen int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}, 3},
	OpAdd:      {"OpAdd", []int{}, 1},
	OpSub:      {"OpSub", []int{}, 1},
	OpMul:      {"OpMul", []int{}, 1},
	OpDiv:      {"OpDiv", []int{}, 1},
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
			binary.BigEndian.PutUint16(bytes[offset:], uint16(operand))
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
