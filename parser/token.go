package parser

type TokenType = *TreeType[string]

// possible TokenType, string is used for debug
var (
	TkIdentifier TokenType = defineType("Identifier")
	TkKeyword    TokenType = defineType("Keyword")
	TkNumber     TokenType = defineType("Number")
	TkString     TokenType = defineType("String")
	TkEof        TokenType = defineType("Eof")

	TkOperator TokenType = defineType("Operator")
	TkAdd      TokenType = defineSubType("Add", TkOperator)
	TkSub      TokenType = defineSubType("Sub", TkOperator)
	TkMul      TokenType = defineSubType("Mul", TkOperator)
	TkDiv      TokenType = defineSubType("Div", TkOperator)
	TkAssign   TokenType = defineSubType("Assign", TkOperator)
	TkPlusEq   TokenType = defineSubType("PlusEq", TkOperator)
	TkSubEq    TokenType = defineSubType("SubEq", TkOperator)
	TkMulEq    TokenType = defineSubType("MulEq", TkOperator)
	TkDivEq    TokenType = defineSubType("DivEq", TkOperator)

	TkDelimiter TokenType = defineType("Delimiter")
	TkComma     TokenType = defineSubType("Comma", TkDelimiter)
	TkSemicolon TokenType = defineSubType("Semicolon", TkDelimiter)
	TkLParen    TokenType = defineSubType("LParen", TkDelimiter)
	TkRParen    TokenType = defineSubType("RParen", TkDelimiter)
	TkLBrace    TokenType = defineSubType("LBrace", TkDelimiter)
	TkRBrace    TokenType = defineSubType("RBrace", TkDelimiter)
	TkLBracket  TokenType = defineSubType("LBracket", TkDelimiter)
	TkRBracket  TokenType = defineSubType("RBracket", TkDelimiter)
)

type Token struct {
	Type    TokenType
	Literal []byte
}
