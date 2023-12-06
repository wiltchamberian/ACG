package parser

import (
	"errors"
	"os"
)

type Lexer struct {
	content       []byte
	rover         Int
	length        Int
	currentLineNo Int
}

func (s *Lexer) SkipWhiteSpaces() Int {
	var count Int = 0
	for s.rover < s.length && s.content[s.rover] == ' ' {
		count += 1
		s.rover += 1
	}
	return count
}

func (s *Lexer) SkipLineBreaks() Int {
	var count Int = 0
	for s.rover < s.length && s.content[s.rover] == '\n' {
		count += 1
		s.rover += 1
	}
	return count
}

func (s *Lexer) SkipWhiteSpacesAndLineBreaks() Int {
	var count Int = 0
	for s.rover < s.length && (s.content[s.rover] == '\n' || s.content[s.rover] == ' ') {
		if s.content[s.rover] == '\n' {
			s.currentLineNo += 1
		}
		count += 1
		s.rover += 1
	}
	return count
}

// skip whitespaces, linebreaks, tabs and so on if possible
func (s *Lexer) SkipAllUnUsed() Int {
	var count Int = 0
	for s.rover < s.length && (s.content[s.rover] == '\n' || s.content[s.rover] == ' ' || s.content[s.rover] == '\t' || s.content[s.rover] == '\r') {
		if s.content[s.rover] == '\n' {
			s.currentLineNo += 1
		}
		count += 1
		s.rover += 1
	}
	return count
}

func (s *Lexer) CheckEnd() bool {
	return s.rover >= s.length
}

func (s *Lexer) GetCharacterType(ch byte) int {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
		return Letter
	} else if ch >= '0' && ch <= '9' {
		return Number
	} else if ch == '_' {
		return UnderLine
	} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '<' || ch == '>' || ch == '=' || ch == '|' {
		return Operator
	} else if ch == ';' || ch == '(' || ch == ')' || ch == '[' || ch == ']' || ch == '{' || ch == '}' || ch == ':' {
		return Delimiter
	} else if ch == '"' || ch == '\'' {
		return String
	}
	return Undefined
}

func (s *Lexer) parseNumber() (Token, error) {
	var token Token
	start := s.rover
	for s.rover < s.length {
		typ := s.GetCharacterType(s.content[s.rover])
		if typ == Number {
			s.rover += 1
			continue
		} else {
			break
		}
	}
	token.Type = TkNumber
	token.Literal = s.content[start:s.rover]
	return token, nil
}

// parse identifier and parse terminator
func (s *Lexer) parseIdentifier() (Token, error) {
	var token Token
	var ec error

	start := s.rover
	s.rover += 1
	for s.rover < s.length {
		typ := s.GetCharacterType(s.content[s.rover])
		if typ == Number || typ == UnderLine || typ == Letter {
			s.rover += 1
			continue
		} else {
			break
		}
	}

	token.Type = TkIdentifier
	token.Literal = s.content[start:s.rover]

	return token, ec
}

func (s *Lexer) parseOperator() (Token, error) {
	var token Token
	var ec error

	ch := s.content[s.rover]
	start := s.rover
	s.rover += 1
	for s.rover < s.length {
		if ch == '+' && s.content[s.rover] == '=' {
			token.Type = TkPlusEq
		} else if ch == '-' && s.content[s.rover] == '=' {
			token.Type = TkSubEq
		} else if ch == '*' && s.content[s.rover] == '=' {
			token.Type = TkMulEq
		} else if ch == '/' && s.content[s.rover] == '=' {
			token.Type = TkDivEq
		} else {
			if ch == '+' {
				token.Type = TkAdd
			} else if ch == '-' {
				token.Type = TkSub
			} else if ch == '*' {
				token.Type = TkMul
			} else if ch == '/' {
				token.Type = TkDiv
			} else if ch == '=' {
				token.Type = TkAssign
			} else if ch == '|' {
				token.Type = TkBitwiseOr
			}
			break
		}
		s.rover += 1
		break

	}
	token.Literal = s.content[start:s.rover]

	return token, ec
}

func (s *Lexer) parseDelimiter() (Token, error) {
	var token Token
	var ec error

	ch := s.content[s.rover]
	start := s.rover
	s.rover += 1

	if ch == ';' {
		token.Type = TkSemicolon
	} else if ch == '(' {
		token.Type = TkLParen
	} else if ch == ')' {
		token.Type = TkRParen
	} else if ch == '{' {
		token.Type = TkLBrace
	} else if ch == '}' {
		token.Type = TkRBrace
	} else if ch == '[' {
		token.Type = TkLBracket
	} else if ch == ']' {
		token.Type = TkRBracket
	} else if ch == ':' {
		token.Type = TkColon
	}

	token.Literal = s.content[start:s.rover]

	return token, ec
}

// TODO：not support linebreak
func (s *Lexer) parseString(ch byte) (Token, error) {
	var token Token
	token.Type = TkString
	var ec error

	start := s.rover
	s.rover += 1

	for s.rover < s.length && s.content[s.rover] != ch {
		s.rover++
	}
	if s.rover >= s.length {
		return token, errors.New("unmatch quotes")
	}
	s.rover += 1

	token.Literal = s.content[start:s.rover]

	return token, ec
}

func (s *Lexer) NextToken() (Token, error) {
	var token Token
	var ec error
	s.SkipAllUnUsed()
	if s.CheckEnd() {
		return token, errors.New("file end")
	}
	ch := s.content[s.rover]
	typ := s.GetCharacterType(ch)
	switch typ {
	case Number:
		{
			token, ec = s.parseNumber()
		}
	case UnderLine, Letter:
		{
			token, ec = s.parseIdentifier()
		}
	case Operator:
		{
			token, ec = s.parseOperator()
		}
	case Delimiter:
		{
			token, ec = s.parseDelimiter()
		}
	case String:
		{
			token, ec = s.parseString(ch)
		}
	default:
		ec = errors.New("undefined character")
	}
	return token, ec
}

func (s *Lexer) ReadString(content string) error {
	var err error
	if content == "" {
		return errors.New("empty content")
	}
	s.content = []byte(content) //FIX ME:assume copy
	s.length = len(content)
	return err
}

func (s *Lexer) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	buffer := make([]byte, fileSize)
	n, err := file.Read(buffer)
	if err != nil {
		return err
	}
	s.content = buffer
	s.length = n
	return nil
}

func (s *Lexer) TokenStream() ([]Token, error) {
	var tokens []Token
	if s.length <= 0 || s.content == nil {
		return tokens, errors.New("empty content")
	}
	var token Token
	var err error
	for {
		token, err = s.NextToken()
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens, err
}
