package parser

import (
	"errors"
	"reflect"
)

type BasicParser struct {
	lexer  Lexer
	tokens []Token
	index  Int
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
	s.index += 1
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
	if token.Type == typ {
		s.index += 1
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
			s.index += 1
			return token, nil
		}
	}
	if reflect.DeepEqual(token.Literal, content) {
		s.index += 1
		return token, nil
	}
	return nil, errors.New("not match")
}

func (s *BasicParser) parse() error {
	return nil
}
