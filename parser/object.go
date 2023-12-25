package parser

import "strconv"

type NkObject interface {
	Type() string
	ToString() string
}

type NkBool struct {
	Value bool
}

func (s *NkBool) Type() string {
	return "Bool"
}

func (s *NkBool) ToString() string {
	if s.Value {
		return "true"
	} else {
		return "false"
	}
}

type NkInteger struct {
	Value int
}

func (s *NkInteger) Type() string {
	return "Integer"
}

func (s *NkInteger) ToString() string {
	return strconv.Itoa(s.Value)
}

type NkFloat struct {
	Value float32
}

func (s *NkFloat) ToString() string {
	return ""
}

func (s *NkFloat) Type() string {
	return "Float"
}

type NkDouble struct {
	Value float64
}

type NkString struct {
	Value string
}

func (s *NkString) Type() string {
	return "String"
}

func (s *NkString) ToString() string {
	return s.Value
}

func Add(left NkObject, right NkObject) NkObject {
	if left.Type() == "Integer" && right.Type() == "Integer" {
		return &NkInteger{Value: (left.(*NkInteger).Value + right.(*NkInteger).Value)}
	} else if left.Type() == "Float" && right.Type() == "Float" {
		return &NkFloat{Value: (left.(*NkFloat).Value + right.(*NkFloat).Value)}
	} else if left.Type() == "String" && right.Type() == "String" {
		return &NkString{Value: (left.(*NkString).Value + right.(*NkString).Value)}
	}
	return nil
}

func Sub(left NkObject, right NkObject) NkObject {
	if left.Type() == "Integer" && right.Type() == "Integer" {
		return &NkInteger{Value: (left.(*NkInteger).Value - right.(*NkInteger).Value)}
	} else if left.Type() == "Float" && right.Type() == "Float" {
		return &NkFloat{Value: (left.(*NkFloat).Value - right.(*NkFloat).Value)}
	}
	return nil
}

func Mul(left NkObject, right NkObject) NkObject {
	if left.Type() == "Integer" && right.Type() == "Integer" {
		return &NkInteger{Value: (left.(*NkInteger).Value * right.(*NkInteger).Value)}
	} else if left.Type() == "Float" && right.Type() == "Float" {
		return &NkFloat{Value: (left.(*NkFloat).Value * right.(*NkFloat).Value)}
	}
	return nil
}

func Div(left NkObject, right NkObject) NkObject {
	if left.Type() == "Integer" && right.Type() == "Integer" {
		return &NkInteger{Value: (left.(*NkInteger).Value / right.(*NkInteger).Value)}
	} else if left.Type() == "Float" && right.Type() == "Float" {
		return &NkFloat{Value: (left.(*NkFloat).Value / right.(*NkFloat).Value)}
	}
	return nil
}
