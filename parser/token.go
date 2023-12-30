package parser

var parentMap = make(map[string]string)

func isType(a string, b string) bool {
	if a == b {
		return true
	}
	for parentMap[a] != "" && parentMap[a] != b {
		a = parentMap[a]
	}
	if parentMap[a] == b {
		return true
	}
	return false
}

//type UsedType = TreeType[string]
// type TokenType = *UsedType

type TokenType = string

func initType(cur string, parent string) TokenType {
	parentMap[cur] = parent
	return cur
}

// possible TokenType, string is used for debug
var (
	TkIdentifier    TokenType = initType("Identifier", "")
	TkTerminator    TokenType = initType("Terminator", TkIdentifier)
	TkNonTerminator TokenType = initType("NonTerminator", TkIdentifier)

	TkEmpty TokenType = initType("Empty", "") //空终结符

	TkKeyword TokenType = initType("Keyword", "")
	TkLet     TokenType = initType("let", TkKeyword)
	TkVar     TokenType = initType("var", TkKeyword)
	TkIf      TokenType = initType("if", TkKeyword)
	TkElse    TokenType = initType("else", TkKeyword)
	TkFor     TokenType = initType("for", TkKeyword)

	TkNumber TokenType = initType("Number", "")
	TkString TokenType = initType("String", "")
	TkEof    TokenType = initType("Eof", "")

	TkOperator   TokenType = initType("Operator", "")
	TkAdd        TokenType = initType("+", TkOperator)
	TkSub        TokenType = initType("-", TkOperator)
	TkMul        TokenType = initType("*", TkOperator)
	TkDiv        TokenType = initType("/", TkOperator)
	TkAssign     TokenType = initType("=", TkOperator)
	TkBitwiseOr  TokenType = initType("|", TkOperator)
	TkBitwiseAnd TokenType = initType("&", TkOperator)
	TkPlusEq     TokenType = initType("+=", TkOperator)
	TkSubEq      TokenType = initType("-=", TkOperator)
	TkMulEq      TokenType = initType("*=", TkOperator)
	TkDivEq      TokenType = initType("/=", TkOperator)
	TkLess       TokenType = initType("<", TkOperator)
	TkGreater    TokenType = initType(">", TkOperator)
	TkLessEq     TokenType = initType("<=", TkOperator)
	TkGreaterEq  TokenType = initType(">=", TkOperator)
	TkEqual      TokenType = initType("==", TkOperator)
	TkNotEq      TokenType = initType("!=", TkOperator)
	TkLShift     TokenType = initType("<<", TkOperator)
	TkRShift     TokenType = initType(">>", TkOperator)
	TkOr         TokenType = initType("||", TkOperator)
	TkAnd        TokenType = initType("&&", TkOperator)

	TkDelimiter TokenType = initType("Delimiter", "")
	TkComma     TokenType = initType(",", TkDelimiter)
	TkSemicolon TokenType = initType(";", TkDelimiter)
	TkLParen    TokenType = initType("(", TkDelimiter)
	TkRParen    TokenType = initType(")", TkDelimiter)
	TkLBrace    TokenType = initType("{", TkDelimiter)
	TkRBrace    TokenType = initType("}", TkDelimiter)
	TkLBracket  TokenType = initType("[", TkDelimiter)
	TkRBracket  TokenType = initType("]", TkDelimiter)
	TkColon     TokenType = initType(":", TkDelimiter)

	TkAction TokenType = initType("Action", "")
)

type Token struct {
	Type    TokenType
	Literal []byte
}

var EmptyToken = &Token{Type: TkEmpty}

func (s *Token) GetLiteral() string {
	return string(s.Literal)
}

func (s *Token) SetName(name string) {
}

func (s *Token) GetName() string {
	return string(s.Type)
}

func (s *Token) GetChildren() []INode {
	return []INode{}
}

func (s *Token) GetParent() INode {
	return nil
}

func (s *Token) Select() int {
	return -1
}

func (s *Token) IsTerminal() bool {
	return true
}
