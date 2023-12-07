package parser

type UsedType = TreeType[string]
type TokenType = *UsedType

// possible TokenType, string is used for debug
var (
	TkIdentifier    TokenType = &UsedType{"Identifier", nil}
	TkTerminator    TokenType = &UsedType{"Terminator", TkIdentifier}
	TkNonTerminator TokenType = &UsedType{"NonTerminator", TkIdentifier}

	TkKeyword TokenType = &UsedType{"Keyword", nil}
	TkNumber  TokenType = &UsedType{"Number", nil}
	TkString  TokenType = &UsedType{"String", nil}
	TkEof     TokenType = &UsedType{"Eof", nil}

	TkOperator  TokenType = &UsedType{"Operator", nil}
	TkAdd       TokenType = &UsedType{"Add", TkOperator}
	TkSub       TokenType = &UsedType{"Sub", TkOperator}
	TkMul       TokenType = &UsedType{"Mul", TkOperator}
	TkDiv       TokenType = &UsedType{"Div", TkOperator}
	TkAssign    TokenType = &UsedType{"Assign", TkOperator}
	TkBitwiseOr TokenType = &UsedType{"BitwiseOr", TkOperator}
	TkPlusEq    TokenType = &UsedType{"PlusEq", TkOperator}
	TkSubEq     TokenType = &UsedType{"SubEq", TkOperator}
	TkMulEq     TokenType = &UsedType{"MulEq", TkOperator}
	TkDivEq     TokenType = &UsedType{"DivEq", TkOperator}

	TkDelimiter TokenType = &UsedType{"Delimiter", nil}
	TkComma     TokenType = &UsedType{"Comma", TkDelimiter}
	TkSemicolon TokenType = &UsedType{"Semicolon", TkDelimiter}
	TkLParen    TokenType = &UsedType{"LParen", TkDelimiter}
	TkRParen    TokenType = &UsedType{"RParen", TkDelimiter}
	TkLBrace    TokenType = &UsedType{"LBrace", TkDelimiter}
	TkRBrace    TokenType = &UsedType{"RBrace", TkDelimiter}
	TkLBracket  TokenType = &UsedType{"LBracket", TkDelimiter}
	TkRBracket  TokenType = &UsedType{"RBracket", TkDelimiter}
	TkColon     TokenType = &UsedType{"Colon", TkDelimiter}
)

type Token struct {
	Type    TokenType
	Literal []byte
}

func (s *Token) GetName() string {
	return string(s.Literal)
}

func (s *Token) GetChildren() []INode {
	return []INode{}
}
