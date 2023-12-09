package main

import (
	. "ACG/parser"
	"fmt"
	"slices"
)

func tt() {
	var arr = []int{1, 2, 3}
	var arr2 = []int{4, 5}
	arr = append(arr, arr2...)
	slices.Reverse(arr)
	for _, ele := range arr {
		fmt.Printf("%d,", ele)
	}
}

func test() {
	var abc []byte = []byte("1234")
	q := abc[0:2]
	abc[0] = 'x'

	fmt.Println(string(abc))
	fmt.Println(string(q))
}

func testLexer() {
	var content string = `x = 13;
	y = 7 + 4;
	z = x * 3;
	q = 'abc123'`

	var lexer Lexer
	err := lexer.ReadString(content)
	if err != nil {
		panic("readstring error")
	}

	tokens, err := lexer.TokenStream()

	for _, token := range tokens {
		fmt.Println(string(token.Literal))
	}
}

func testEbnfParser() ([]Rule, error) {
	path := "./nika_simple.gram"
	var parser EBNFParser
	parser.ReadFile(path)
	parser.TokenStream()
	tokens := parser.GetTokens()
	for _, token := range tokens {
		fmt.Println(string(token.Literal))
	}
	rules := parser.Grammar()
	for _, rule := range rules {
		fmt.Println(rule.Name)
	}
	return rules, nil
}

func testGenerator() {
	rules, _ := testEbnfParser()
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	name := "NikaParser"
	generator.Generate_rparser(name, rules)

}
