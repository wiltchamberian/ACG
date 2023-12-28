package parser

import (
	"errors"
	"reflect"
)

// control left or right recursion parse

type BasicParser struct {
	lexer   Lexer
	tokens  []Token
	index   Int
	control ParserControllor
}

func NewBasicParserL() BasicParser {
	var v BasicParser
	v.control = &ParserControllorL{}
	return v
}

func NewBasicParserR() BasicParser {
	var v BasicParser
	v.control = &ParserControllorR{}
	return v
}

func (s *BasicParser) Clear() {
	s.lexer.Reset()
	s.tokens = nil
	s.index = 0
}

func (s *BasicParser) ReadFile(path string) error {
	return s.lexer.ReadFile(path)
}

func (s *BasicParser) ReadString(path string) error {
	return s.lexer.ReadString(path)
}

func (s *BasicParser) TokenStream() error {
	var err error
	s.tokens, err = s.lexer.TokenStream()
	s.control.ResetIndex(s)
	return err
}

func (s *BasicParser) Tokenize(str string) error {
	err := s.ReadString(str)
	if err != nil {
		return err
	}
	err = s.TokenStream()
	return err
}

func (s *BasicParser) GetTokens() []Token {
	return s.tokens
}

func (s *BasicParser) PeekToken() (*Token, error) {
	if s.index >= len(s.tokens) || s.index < 0 {
		return nil, errors.New("overflow")
	}
	return &s.tokens[s.index], nil
}

func (s *BasicParser) GetToken() (*Token, error) {
	if s.index >= len(s.tokens) || s.index < 0 {
		return nil, errors.New("overflow")
	}
	s.control.AdvanceIndex(s)
	return &s.tokens[s.index-1], nil
}

func (s *BasicParser) Mark() Int {
	return s.index
}

func (s *BasicParser) Reset(pos Int) {
	s.index = pos
}

func (s *BasicParser) Expect(typ TokenType) (*Token, error) {
	token, err := s.PeekToken()
	if err != nil {
		return nil, err
	}
	if isType(token.Type, typ) {
		s.control.AdvanceIndex(s)
		return token, err
	}
	return nil, errors.New("not match")
}

func (s *BasicParser) ExpectValue(content interface{}) (*Token, error) {
	token, err := s.PeekToken()
	if err != nil {
		return nil, err
	}
	if value, ok := content.(string); ok == true {
		if string(token.Literal) == value {
			s.control.AdvanceIndex(s)
			return token, nil
		}
	}
	if reflect.DeepEqual(token.Literal, content) {
		s.control.AdvanceIndex(s)
		return token, nil
	}
	return nil, errors.New("not match")
}

func (s *BasicParser) parse() error {
	return nil
}

func (s *BasicParser) ExpectW(typ TokenType) Ret {
	a, b := s.Expect(typ)
	return Ret{a, b}
}

func (s *BasicParser) ExpectValueW(content interface{}) Ret {
	a, b := s.ExpectValue(content)
	return Ret{a, b}
}
